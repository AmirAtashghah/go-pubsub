package main

import (
	"go-pubsub/real-example/api"
	boardmaker "go-pubsub/real-example/board"
	"net/http"

	"github.com/labstack/echo/v4"
)

func init() {

	boardmaker.Receiver()
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, welcome to the user registration app!")
	})

	e.POST("/register", api.RegisterUser)

	e.Logger.Fatal(e.Start(":8080"))
}
