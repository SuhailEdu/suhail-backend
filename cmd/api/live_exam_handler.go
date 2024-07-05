package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	_ "golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
)

func (config *Config) getLiveExam(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	isParticipant, err := isExamParticipant(c, config, examId, userId)
	if err != nil {
		return serverError(c, err)
	}
	if !isParticipant {
		return unAuthorizedError(c, "unauthorized user")
	}

	exam, err := config.db.GetExamById(c.Request().Context(), examId)

	if err != nil {
		return serverError(c, err)
	}

	questions, err := config.db.GetExamQuestions(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	isIpAllowed := true
	if exam.IpRangeStart.String != "" && exam.IpRangeEnd.String != "" {
		start := net.ParseIP(exam.IpRangeStart.String)
		end := net.ParseIP(exam.IpRangeEnd.String)

		userIP := net.ParseIP(c.RealIP())
		fmt.Println(start.To4())
		fmt.Println(end.To4())
		fmt.Println(userIP.To4())

		if start.To4() != nil && end.To4() != nil && userIP.To4() != nil {
			fmt.Println("here")

			if bytes.Compare(start, userIP) > 0 || bytes.Compare(end, userIP) < 0 {
				isIpAllowed = false
			}

		}
	}

	return dataResponse(c, types.SerializeGetLiveExam(exam, questions, isIpAllowed))

}

func (config *Config) getLiveExamParticipantsForManager(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	isAuther, _ := isExamAuthor(c, config, examId)

	if !isAuther {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	participants, err := config.db.GetLiveExamParticipants(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	ps := getExamParticipantsConnectionStatus(config, examId)

	return dataResponse(c, types.SerializeGetLiveExamParticipants(participants, ps))
}
func getExamParticipantsConnectionStatus(config *Config, examId uuid.UUID) []uuid.UUID {

	var ps []uuid.UUID

	sessions, err := config.melody.Sessions()
	if err != nil {
		return []uuid.UUID{}
	}

	for _, s := range sessions {
		sessionExamId, ok := s.Keys["examId"].(uuid.UUID)
		if !ok {
			continue
		}

		if examId != sessionExamId {
			continue
		}

		isAuthor, ok := s.Keys["isAuthor"].(bool)
		if !ok {
			continue
		}

		if isAuthor {
			continue
		}

		participantId, ok := s.Keys["userId"].(uuid.UUID)

		if !ok {
			continue
		}

		ps = append(ps, participantId)
	}

	return ps

}

func (config *Config) getLiveExamQuestionsForManager(c echo.Context) error {
	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	questions, err := config.db.GetLiveExamQuestionForManager(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeGetLiveExamQuestionsForManager(questions))
}

func (config *Config) updateExamStatus(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	var body struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.Bind(&body); err != nil {
		return badRequestError(c, errors.New("invalid status"))
	}

	if body.Status != "paused" && body.Status != "finished" && body.Status != "live" {
		return badRequestError(c, errors.New("invalid status"))
	}

	isAuther, _ := isExamAuthor(c, config, examId)
	if !isAuther {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	err = config.db.UpdateExamLiveStatus(c.Request().Context(), schema.UpdateExamLiveStatusParams{
		LiveStatus: pgtype.Text{String: body.Status, Valid: true},
		ID:         examId,
	})

	if err != nil {
		return serverError(c, err)
	}

	broadcastLiveExamStatusUpdate(c, config, body.Status)

	return c.NoContent(http.StatusOK)
}

func (config *Config) storeAnswer(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	var body struct {
		QuestionId uuid.UUID `json:"question_id"`
		Answer     string    `json:"answer"`
	}

	if err := c.Bind(&body); err != nil {
		return badRequestError(c, errors.New("invalid body"))
	}

	if _, err = body.QuestionId.Value(); err != nil {
		return badRequestError(c, errors.New("invalid question id"))
	}

	if body.Answer == "" {
		return badRequestError(c, errors.New("invalid answer"))
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	isParticipant, _ := isExamParticipant(c, config, examId, userId)
	if !isParticipant {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	//questionUUID, err := uuid.FromBytes([]byte(body.QuestionId))
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println(body.QuestionId)
	//	return badRequestError(c, errors.New("invalid question id"))
	//}

	fmt.Println(body.QuestionId)

	isQuestionExits, err := config.db.CheckQuestionExits(c.Request().Context(), schema.CheckQuestionExitsParams{
		ID:     body.QuestionId,
		ExamID: examId,
	})

	if !isQuestionExits {
		return badRequestError(c, errors.New("invalid question id"))
	}

	err = config.db.UpdateAnswer(c.Request().Context(), schema.UpdateAnswerParams{
		QuestionID: body.QuestionId,
		UserID:     authenticatedUser.ID,
		Answer:     body.Answer,
	})

	if err != nil {
		return serverError(c, err)
	}

	return c.NoContent(http.StatusOK)
}
func isExamParticipant(c echo.Context, config *Config, examId uuid.UUID, participantId uuid.UUID) (bool, error) {

	uuidValue := pgtype.UUID{
		Bytes: participantId,
		Valid: true,
	}
	isParticipant, err := config.db.CheckParticipant(c.Request().Context(), schema.CheckParticipantParams{
		ExamID: examId,
		UserID: uuidValue,
	})

	if err != nil {
		return false, err
	}

	return isParticipant, nil

}

func broadcastLiveExamStatusUpdate(c echo.Context, config *Config, status string) {

	statusEvent := map[string]interface{}{
		"type": "LIVE_EXAM_STATUS_UPDATED",
		"payload": map[string]interface{}{
			"status": status,
		},
	}

	jsonEvent, _ := json.Marshal(statusEvent)
	//err := config.melody.Broadcast(jsonEvent)
	//if err != nil {
	//	config.logger.Println(err)
	//	return
	//}

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return
	}

	_ = config.melody.BroadcastFilter(jsonEvent, func(s *melody.Session) bool {
		fmt.Println(s.Request.Header)

		authHeader := s.Request.Header.Get("Sec-WebSocket-Protocol")
		if authHeader == "" {
			fmt.Println("auth header is empty")
			return false

		}

		hash := sha256.Sum256([]byte(authHeader))

		isParticipant, err := config.db.CheckParticipantByToken(s.Request.Context(), schema.CheckParticipantByTokenParams{
			ExamID: examId,
			Hash:   hash[:],
		})

		if err != nil {
			fmt.Println(err)
			return false
		}

		return isParticipant
	})

}
