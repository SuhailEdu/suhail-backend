package main

import (
	"errors"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
)

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

	return dataResponse(c, types.SerializeGetLiveExamParticipants(participants))
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

	status := c.FormValue("status")

	if status != "paused" && status != "finished" && status != "live" {
		return badRequestError(c, errors.New("invalid status"))
	}

	isAuther, _ := isExamAuthor(c, config, examId)
	if !isAuther {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	err = config.db.UpdateExamLiveStatus(c.Request().Context(), schema.UpdateExamLiveStatusParams{
		LiveStatus: pgtype.Text{String: status, Valid: true},
		ID:         examId,
	})

	if err != nil {
		return serverError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (config *Config) storeAnswer(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	questionId := c.FormValue("questionId")
	answer := c.FormValue("answer")

	if questionId == "" {
		return badRequestError(c, errors.New("invalid question id"))
	}

	if answer == "" {
		return badRequestError(c, errors.New("invalid answer"))
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	isParticipant, _ := isExamParticipant(c, config, examId, userId)
	if !isParticipant {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	questionUUID, err := uuid.FromBytes([]byte(questionId))
	if err != nil {
		return badRequestError(c, errors.New("invalid question id"))
	}

	isQuestionExits, err := config.db.CheckQuestionExits(c.Request().Context(), schema.CheckQuestionExitsParams{
		ID:     questionUUID,
		ExamID: examId,
	})

	if !isQuestionExits {
		return badRequestError(c, errors.New("invalid question id"))
	}

	err = config.db.UpdateAnswer(c.Request().Context(), schema.UpdateAnswerParams{
		QuestionID: questionUUID,
		UserID:     authenticatedUser.ID,
		Answer:     answer,
	})

	if err != nil {
		return serverError(c, err)
	}

	return c.NoContent(http.StatusOK)
}
func isExamParticipant(c echo.Context, config *Config, examId uuid.UUID, participantId uuid.UUID) (bool, error) {

	uuidValue := pgtype.UUID{
		Bytes: participantId,
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
