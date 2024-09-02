package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type MessageStorage struct {
	messages map[int64]bool
	lock     sync.Mutex
}

func (ms *MessageStorage) addMessage(messageNumber int64) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.messages[messageNumber] = true
}

func (ms *MessageStorage) getMessages() []int64 {
	messages := []int64{}

	for key := range ms.messages {
		messages = append(messages, key)
	}

	return messages
}

func main() {
	node := maelstrom.NewNode()

	ms := MessageStorage{
		messages: make(map[int64]bool),
		lock:     sync.Mutex{},
	}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgInt, ok := body["message"].(float64)
		if !ok {
			return fmt.Errorf("Message not an int")
		}

		ms.addMessage(int64(msgInt))

		delete(body, "message")
		body["type"] = "broadcast_ok"

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "read_ok"
		body["messages"] = ms.getMessages()

		return node.Reply(msg, body)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "topology_ok"
		delete(body, "topology")

		return node.Reply(msg, body)

	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
