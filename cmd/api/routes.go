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

	homeGroup.POST("/exams/create", config.createExam)
	homeGroup.GET("/exams", config.getExamsList)
	homeGroup.GET("/others_exams", config.getParticipatedExams)
	homeGroup.GET("/exams/:id", config.getSingleExam)
	homeGroup.PATCH("/exams/:id", config.updateExam)

	homeGroup.POST("/exams/:id/questions", config.addQuestionsToExam)
	homeGroup.PATCH("/exams/:id/questions/:questionId", config.updateQuestion)
	homeGroup.DELETE("/exams/:id/questions/:questionId", config.deleteQuestion)

	homeGroup.POST("/exams/:id/invite", config.inviteUsersToExam)

	homeGroup.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
