package api

import (
	"encoding/json"
	"fmt"
	"go-pubsub/message"
	"go-pubsub/pubsub/gochannel"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID       string `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func RegisterUser(c echo.Context) error {
	u := new(User)

	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}

	// Generate a unique ID for the user
	u.ID = uuid.New().String()

	go publisher("user.register", *u)

	// Perform user registration logic here, e.g., store the user in a database.

	// For now, let's just print the user information and board information.
	fmt.Printf("User registered: %+v\n", u)

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func publisher(topic string, u User) {

	userJSON, err := json.Marshal(u)
	if err != nil {

		slog.Info("can not marshal")
	}

	msg := message.NewMessage(uuid.New(), userJSON)

	err = gochannel.GoChanPubsub.Publish(topic, msg)
	if err != nil {

		slog.Info("can not publish")
	}

	// Create a default board for the user

}
