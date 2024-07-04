package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
)

type SessionMetadata struct {
	UserId   uuid.UUID `json:"userId"`
	ExamId   uuid.UUID `json:"examId"`
	Token    []byte    `json:"token"`
	IsAuthor bool      `json:"isAuthor"`
}

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

		metadata := map[string]interface{}{
			"examId":   examId,
			"userId":   authenticatedUser.ID,
			"isAuthor": false,
		}

		err = config.melody.HandleRequestWithKeys(c.Response().Writer, c.Request(), metadata)

		return err

	})

	// for manager
	wsGroup.GET("/live/:id/manage", func(c echo.Context) error {

		examId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return badRequestError(c, err)
		}

		authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)

		isAuthor, _ := isExamAuthor(c, config, examId)
		if !isAuthor {
			return unAuthorizedError(c, "unauthorized user")
		}
		metadata := map[string]interface{}{
			"examId":   examId,
			"token":    authenticatedUser.Hash,
			"userId":   authenticatedUser.ID,
			"isAuthor": isAuthor,
		}

		err = config.melody.HandleRequestWithKeys(c.Response().Writer, c.Request(), metadata)
		return err

	})

	config.melody.HandleDisconnect(func(s *melody.Session) {
		fmt.Println("hi")
		fmt.Println(s.Keys["examId"])

		isAuthor, ok := s.Keys["isAuthor"]
		if !ok {
			return

		}
		if !isAuthor.(bool) {
			return
		}

		userId, ok := s.Keys["userId"].(uuid.UUID)
		if !ok {
			return

		}

		examId, ok := s.Keys["examId"].(uuid.UUID)
		if !ok {
			return
		}

		token, ok := s.Keys["token"].([]byte)
		if !ok {
			return
		}

		broadcastParticipantConnection(config, examId, userId, token, false)
	})

	config.melody.HandleConnect(func(s *melody.Session) {
		fmt.Println("hi")
		fmt.Println(s.Keys["examId"])

		isAuthor, ok := s.Keys["isAuthor"]
		if !ok {
			return

		}
		if !isAuthor.(bool) {
			return
		}

		userId, ok := s.Keys["userId"].(uuid.UUID)
		if !ok {
			return

		}

		examId, ok := s.Keys["examId"].(uuid.UUID)
		if !ok {
			return
		}

		token, ok := s.Keys["token"].([]byte)
		if !ok {
			return
		}

		broadcastParticipantConnection(config, examId, userId, token, true)
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

func broadcastParticipantConnection(config *Config, examId uuid.UUID, userId uuid.UUID, token []byte, isConnect bool) {

	var eventType string
	if isConnect {
		eventType = "PARTICIPANT_CONNECTED"
	} else {
		eventType = "PARTICIPANT_DISCONNECTED"
	}
	fmt.Println("broadcasting participant")
	fmt.Println("isConnect:", isConnect)

	statusEvent := map[string]interface{}{
		"type": eventType,
		"payload": map[string]interface{}{
			"participant_id": userId,
		},
	}

	jsonEvent, _ := json.Marshal(statusEvent)

	err := config.melody.BroadcastFilter(jsonEvent, func(s *melody.Session) bool {

		authHeader := s.Request.Header.Get("Sec-WebSocket-Protocol")
		return bytes.Equal([]byte(authHeader), token)

	})

	if err != nil {
		config.logger.Println("broadcast disconnected participant error:", err)
		return
	}

}
