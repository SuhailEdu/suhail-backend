package main

import (
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "golang.org/x/crypto/bcrypt"
)

func (config *Config) getParticipantAnswers(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	participantId, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		return badRequestError(c, err)
	}

	questions, err := config.db.GetParticipantAnswers(c.Request().Context(), schema.GetParticipantAnswersParams{
		ID:     examId,
		UserID: participantId,
	})

	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeGetParticipantAnswers(questions))

}
