package board

import (
	"context"
	"encoding/json"
	"fmt"
	"go-pubsub/example/user"
	"go-pubsub/message"
	"go-pubsub/pubsub/gochannel"
	"log/slog"
)

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

		u := new(user.User)

		err := json.Unmarshal(msg.Payload, u)
		if err != nil {
			slog.Info("Error:", "err :", err)

		}

		fmt.Printf("received message: %s, payload: %s\n", msg.Id.String(), u)

	}
}
