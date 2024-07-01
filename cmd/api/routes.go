package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func registerApiRoutes(e *echo.Echo, config *Config) {

	authGroup := e.Group("auth")

	authGroup.POST("/register", config.registerUser)
	authGroup.POST("/login", config.loginUser)

	homeGroup := e.Group("/home")
	homeGroup.Use(config.checkAuthToken)

	homeGroup.POST("/logout", config.logout)

	homeGroup.GET("/exams", config.getExamsList)
	homeGroup.POST("/exams/create", config.createExam)

	homeGroup.POST("/exams/create/check-title", config.checkExamTitle)
	homeGroup.POST("/exams/create/check-questions", config.checkExamQuestions)

	homeGroup.GET("/others_exams", config.getParticipatedExams)
	homeGroup.GET("/exams/:id", config.getSingleExam)
	homeGroup.PATCH("/exams/:id", config.updateExam)
	homeGroup.DELETE("/exams/:id", config.updateExam)

	homeGroup.GET("/exams/:id/participants", config.getExamParticipants)
	homeGroup.POST("/exams/:id/invite", config.inviteUsersToExam)
	homeGroup.DELETE("/exams/:id/removeParticipants", config.removeParticipants)

	homeGroup.POST("/exams/:id/questions", config.addQuestionsToExam)
	homeGroup.PATCH("/exams/:id/questions/:questionId", config.updateQuestion)
	homeGroup.DELETE("/exams/:id/questions/:questionId", config.deleteQuestion)

	homeGroup.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
