package main

import (
	"errors"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "golang.org/x/crypto/bcrypt"
)

func (config *Config) getLiveExamParticipantsForManager(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	_, err = config.db.GetUserExams(c.Request().Context(), userId)
	if err != nil {
		return serverError(c, err)
	}

	isAuther, _ := isExamAuthor(c, config, examId)

	if !isAuther {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	participants, err := config.db.GetExamParticipants(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeGetExamParticipants(participants))
}

func (config *Config) getLiveExamQuestionsForManager(c echo.Context) error {
	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	questions, err := config.db.GetExamQuestions(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeGetExamQuestions(questions))
}
