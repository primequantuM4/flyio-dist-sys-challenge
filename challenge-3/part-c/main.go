package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ClusterNode struct {
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

func persistBroadcastMessage(srcNode *maelstrom.Node,
	destNode string,
	message,
	msg_id float64,
) {
	body := map[string]any{
		"msg_id":  msg_id,
		"type":    "broadcast",
		"message": message,
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if _, err := srcNode.SyncRPC(ctx, destNode, body); err != nil {
			continue
		}

		return
	}
}

func main() {
	node := maelstrom.NewNode()

	ms := MessageStorage{
		messages: make(map[int64]bool),
		lock:     sync.Mutex{},
	}

	cn := ClusterNode{
		neighborNode: make(map[string]bool),
		lock:         sync.Mutex{},
	}

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology := body["topology"].(map[string]interface{})
		nodeNeighbors := topology[node.ID()].([]interface{})

		for _, nextNode := range nodeNeighbors {
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
		msg_id := body["msg_id"].(float64)

		if !ok {
			return fmt.Errorf("Message not a number")
		}

		for _, nextNode := range cn.getNeighbors() {
			if nextNode == msg.Src {
				continue
			}

			go persistBroadcastMessage(node, nextNode, msgInt, msg_id)
		}

		body["type"] = "broadcast_ok"
		delete(body, "message")

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
