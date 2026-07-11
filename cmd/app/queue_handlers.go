package main

import (
	"fmt"
	"time"

	"github.com/alexandergu/queque/internal/queue"
)

var QueueHandlers = map[string]queue.Handler{
	"test1": func(bytes []byte) ([]byte, error) {
		fmt.Println("first process")
		time.Sleep(1 * time.Second)

		return []byte{}, nil
	},

	"test2": func(bytes []byte) ([]byte, error) {
		fmt.Println("second process")
		time.Sleep(2 * time.Second)

		return []byte{}, nil
	},

	"test3": func(bytes []byte) ([]byte, error) {
		fmt.Println("third process")
		time.Sleep(1 * time.Second)

		return []byte{}, nil
	},
}
