package main

import (
	"encoding/json"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	"github.com/SuhailEdu/suhail-backend/internal/validations"
	_ "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
)

func (config *Config) getExamsList(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID
	exams, err := config.db.GetUserExams(c.Request().Context(), userId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeExams(exams))

}

func (config *Config) getParticipatedExams(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID
	exams, err := config.db.GetParticipatedExams(c.Request().Context(), userId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeParticipatedExams(exams))

}

func (config *Config) getSingleExam(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID
	examId := pgtype.UUID{Bytes: [16]byte([]byte(c.Param("id")))}
	fmt.Println(examId, userId)

	exam, err := config.db.FindMyExam(c.Request().Context(), schema.FindMyExamParams{ID: examId, UserID: userId})

	if err.Error() == "no rows in result set" {
		participatedExam, pErr := config.db.FindMyParticipatedExam(c.Request().Context(), schema.FindMyParticipatedExamParams{ID: examId, UserID: userId})
		if pErr.Error() == "no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{})
		}
		return dataResponse(c, participatedExam)

	}

	return dataResponse(c, types.SerializeSingleExam(exam))

}
func (config *Config) createExam(c echo.Context) error {

	var examSchema types.ExamInput

	rules := govalidator.MapData{
		"exam_title": []string{"required", "min:4", "max:30"},
		"status":     []string{"required", "in:public,private"},
		"questions":  []string{"required"},
	}

	opts := govalidator.Options{
		Request: c.Request(), // request object
		Rules:   rules,       // rules map
		Data:    &examSchema,
	}
	// Create a new validator instance
	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {

		return validationError(c, e)
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)

	examTitleExists, _ := config.db.CheckExamTitleExists(c.Request().Context(), schema.CheckExamTitleExistsParams{
		Title:  examSchema.ExamTitle,
		UserID: authenticatedUser.ID,
	})
	if examTitleExists {
		return validationError(c, map[string]interface{}{
			"exam_title": "You already have an exam with this title.",
		})

	}

	isCorrect, questionErrors := validations.ValidateQuestions(examSchema.Questions)
	if !isCorrect {
		return validationError(c, map[string]interface{}{
			"questions": questionErrors,
		})
	}

	examParams := schema.CreateExamParams{
		UserID:           authenticatedUser.ID,
		Title:            examSchema.ExamTitle,
		Slug:             pgtype.Text{String: examSchema.ExamTitle},
		IsAccessable:     pgtype.Bool{Bool: true},
		VisibilityStatus: examSchema.Status,
	}

	createdExam, err := config.db.CreateExam(c.Request().Context(), examParams)
	if err != nil {
		return serverError(c, err)
	}

	var questionsParams []schema.CreateExamQuestionsParams

	for _, question := range examSchema.Questions {
		answers, err := json.Marshal(question.Options)
		if err != nil {
			return serverError(c, err)
		}
		questionsParams = append(questionsParams, schema.CreateExamQuestionsParams{
			ExamID:   createdExam.ID,
			Question: question.Title,
			Type:     "options",
			Answers:  answers,
		})
	}

	_, err = config.db.CreateExamQuestions(c.Request().Context(), questionsParams)

	if err != nil {
		return serverError(c, err)
	}

	examId, err := createdExam.ID.UUIDValue()
	if err != nil {
		return serverError(c, err)
	}

	createdQuestions, _ := config.db.GetExamQuestions(c.Request().Context(), examId)

	return dataResponse(c, types.SerializeExamResource(createdExam, createdQuestions))

}
