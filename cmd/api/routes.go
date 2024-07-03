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
	homeGroup.GET("/participating-exams", config.getParticipatedExams)
	homeGroup.POST("/exams/create", config.createExam)

	homeGroup.POST("/exams/create/check-title", config.checkExamTitle)
	homeGroup.POST("/exams/create/check-questions", config.checkExamQuestions)

	homeGroup.GET("/exams/:id", config.getSingleExam)
	homeGroup.PATCH("/exams/:id", config.updateExam)
	homeGroup.DELETE("/exams/:id", config.updateExam)
	homeGroup.GET("/exams/:id/live", config.getLiveExam)

	homeGroup.GET("/exams/:id/live/manage/participants", config.getLiveExamParticipantsForManager)
	homeGroup.GET("/exams/:id/live/manage/questions", config.getLiveExamQuestionsForManager)

	homeGroup.POST("/exams/:id/live/manage/update-status", config.updateExamStatus)
	homeGroup.POST("/exams/:id/live/store-answer", config.storeAnswer)

	homeGroup.GET("/exams/:id/participants", config.getExamParticipants)
	homeGroup.POST("/exams/:id/invite", config.inviteUsersToExam)
	homeGroup.POST("/exams/:id/remove-participants", config.removeParticipants)

	homeGroup.GET("/exams/:id/questions", config.getExamQuestions)
	homeGroup.POST("/exams/:id/questions", config.addQuestionsToExam)
	homeGroup.PATCH("/exams/:id/questions/:questionId", config.updateQuestion)
	homeGroup.DELETE("/exams/:id/questions/:questionId", config.deleteQuestion)

	homeGroup.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
