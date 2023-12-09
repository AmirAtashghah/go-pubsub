package main

import (
	"go-pubsub/example/board"
	"go-pubsub/example/user"
	"time"
)

func main() {

	board.Receiver()

	for i := 0; i < 2; i++ {
		user.Register()

		time.Sleep(time.Second)

	}

}
