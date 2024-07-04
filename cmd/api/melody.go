package main

import (
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func registerMelodyHandlers(e *echo.Echo, config *Config) {

	wsGroup := e.Group("/ws")
	wsGroup.Use(config.checkWsAuthToken)

	// for participants
	wsGroup.GET("/live/:id", func(c echo.Context) error {
		fmt.Println("new connection try")
		examId, err := uuid.Parse(c.Param("id"))

		if err != nil {
			return badRequestError(c, err)
		}

		authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
		userId := authenticatedUser.ID
		fmt.Println("new connection try 2")

		isParticipant, err := isExamParticipant(c, config, examId, userId)
		if err != nil {
			return serverError(c, err)
		}
		if !isParticipant {
			return unAuthorizedError(c, "unauthorized user")
		}
		fmt.Println("new connection try 3")

		err = config.melody.HandleRequest(c.Response().Writer, c.Request())
		return err

	})

	// for manager
	wsGroup.GET("/live/:id/manage", func(c echo.Context) error {

		examId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return badRequestError(c, err)
		}

		//authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
		//userId := authenticatedUser.ID

		isAuthor, _ := isExamAuthor(c, config, examId)
		if !isAuthor {
			return unAuthorizedError(c, "unauthorized user")
		}

		err = config.melody.HandleRequest(c.Response().Writer, c.Request())
		return err

	})

	//config.melody.HandleConnect(func(s *melody.Session) {
	//	fmt.Println("New connection!")
	//})
	//
	//config.melody.HandleMessage(func(s *melody.Session, msg []byte) {
	//	err := config.melody.Broadcast(msg)
	//	if err != nil {
	//		return
	//	}
	//	fmt.Println("broadcast", string(msg))
	//})
	//
	//config.melody.HandleError(func(s *melody.Session, err error) {
	//	fmt.Println("Error occurred:", err)
	//})
}
