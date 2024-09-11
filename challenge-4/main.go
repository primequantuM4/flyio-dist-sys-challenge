package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	node.Handle("init", func(msg maelstrom.Message) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := kv.Write(ctx, node.ID(), 0); err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		sum := 0

		for _, nodeID := range node.NodeIDs() {

			if nodeID == node.ID() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				value, err := kv.ReadInt(ctx, node.ID())
				if err != nil {
					log.Printf("Failed to read %s: %v", nodeID, err)
					continue
				}

				sum += value
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				res, err := node.SyncRPC(ctx, nodeID, map[string]any{"type": "local_read"})

				if err != nil {
					log.Printf("failed to read %s: %v", nodeID, err)
					continue
				}

				var localBody map[string]any

				if err := json.Unmarshal(res.Body, &localBody); err != nil {
					return err
				}

				valInt, ok := localBody["value"].(float64)

				if !ok {
					return fmt.Errorf("Message is not a number")
				}
				sum += int(valInt)
			}
		}

		body["type"] = "read_ok"
		body["value"] = sum

		return node.Reply(msg, body)

	})

	node.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgInt, ok := body["delta"].(float64)

		if !ok {
			return fmt.Errorf("Message is not a number %f", msgInt)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		sum, err := kv.ReadInt(ctx, node.ID())

		if err != nil {
			return fmt.Errorf("The problem is with the %s", node.ID())
		}

		ctx, cancelWrite := context.WithTimeout(context.Background(), time.Second)
		defer cancelWrite()

		kv.Write(ctx, node.ID(), sum+int(msgInt))

		delete(body, "delta")
		body["type"] = "add_ok"

		return node.Reply(msg, body)
	})

	node.Handle("local_read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		value, err := kv.ReadInt(ctx, node.ID())
		if err != nil {
			return err
		}

		body["type"] = "local_read_ok"
		body["value"] = value

		return node.Reply(msg, body)

	})
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
