package nats

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	log.Printf("[NATS]Connected to %s", url)
	return &Publisher{conn: nc}, nil
}

func (p *Publisher) PublishOrderCreated(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("[NATS]JSON marshal failed: %v", err)
		return err
	}

	subject := "order.created"
	log.Printf("[NATS]Publishing: %s → %s", subject, string(bytes))

	return p.conn.Publish(subject, bytes)
}

func (p *Publisher) Close() {
	if p.conn != nil && !p.conn.IsClosed() {
		p.conn.Close()
		log.Println("[NATS]Connection closed")
	}
}
