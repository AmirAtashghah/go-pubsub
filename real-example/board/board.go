package boardmaker

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pubsub/message"
	"go-pubsub/pubsub/gochannel"
	"go-pubsub/real-example/api"
	"log/slog"

	"github.com/google/uuid"
)

type Board struct {
	ID   string   `json:"id"`
	User api.User `json:"user"`
}

func createBoard(user api.User) {
	defaultBoard := Board{
		ID: uuid.New().String(),

		User: user,
	}

	fmt.Printf("Default board created: %+v\n", defaultBoard)
}

func Receiver() {

	chanPubsub := gochannel.GoChanPubsub

	messages, err := chanPubsub.Subscribe(context.Background(), "user.register")
	if err != nil {
		panic(err)
	}

	go process(messages)

}

func process(messages <-chan *message.Message) {
	for msg := range messages {

		u := new(api.User)

		err := json.Unmarshal(msg.Payload, u)
		if err != nil {
			slog.Info("Error:", "err :", err)

		}

		createBoard(*u)

		//fmt.Printf("received message: %s, payload: %s\n", msg.Id.String(), u)

	}
}
