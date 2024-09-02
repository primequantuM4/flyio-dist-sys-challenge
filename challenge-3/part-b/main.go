package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ClusterNode struct {
	currNode     string
	neighborNode map[string]bool
	lock         sync.Mutex
}

func (cn *ClusterNode) addNeighbor(neighbor string) {
	cn.lock.Lock()
	defer cn.lock.Unlock()
	cn.neighborNode[neighbor] = true
}

func (cn *ClusterNode) getNeighbors() []string {
	nodes := []string{}

	for key := range cn.neighborNode {
		nodes = append(nodes, key)
	}

	return nodes
}

type MessageStorage struct {
	messages map[int64]bool
	lock     sync.Mutex
}

func (ms *MessageStorage) addMessage(message int64) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.messages[message] = true
}

func (ms *MessageStorage) getMessages() []int64 {
	messages := []int64{}

	for key := range ms.messages {
		messages = append(messages, key)
	}

	return messages
}

func (ms *MessageStorage) messageExsists(message int64) bool {
	_, ok := ms.messages[message]

	return ok

}
func main() {
	node := maelstrom.NewNode()
	ms := MessageStorage{
		messages: make(map[int64]bool),
		lock:     sync.Mutex{},
	}

	cn := ClusterNode{
		neighborNode: make(map[string]bool),
	}

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology := body["topology"].(map[string]interface{})
		currNodeNeighbors := topology[node.ID()].([]interface{})

		for _, nextNode := range currNodeNeighbors {
			cn.addNeighbor(nextNode.(string))
		}

		body["type"] = "topology_ok"
		delete(body, "topology")

		return node.Reply(msg, body)
	})

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgInt, ok := body["message"].(float64)

		if !ok {
			return fmt.Errorf("Message is not a number")
		}

		for _, nextNode := range cn.getNeighbors() {
			if nextNode == msg.Src {
				continue
			}

			node.Send(nextNode, body)
		}

		delete(body, "message")
		body["type"] = "broadcast_ok"

		if ms.messageExsists(int64(msgInt)) {
			node.Reply(msg, body)
			return nil
		}

		ms.addMessage(int64(msgInt))

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

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}

}
