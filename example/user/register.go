package user

import (
	"encoding/json"
	"go-pubsub/message"
	"go-pubsub/pubsub/gochannel"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id    string
	Email string
}

func Register() {

	u := new(User)

	u.Id = "1111111"
	u.Email = "amir@gmail.com"

	userJSON, err := json.Marshal(u)
	if err != nil {

		slog.Info("can not marshal")
	}

	msg := message.NewMessage(uuid.New(), userJSON)

	err = gochannel.GoChanPubsub.Publish("user.register", msg)
	if err != nil {

		slog.Info("can not publish")
	}

	time.Sleep(2 * time.Second)
}
