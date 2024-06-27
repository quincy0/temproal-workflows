package mq

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	MyTopic   = "workflow-data"
	Partition = 0
)

var (
	kafkaAddress = "localhost:19092"
)

func Produce(data []byte) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, MyTopic, Partition)
	if err != nil {
		log.Fatal("failed to dial:", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.WriteMessages(kafka.Message{Value: data})
	if err != nil {
		log.Fatal("failed to write message:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

var Reader *kafka.Reader

var MsgChan chan string

func init() {
	Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddress},
		Topic:   MyTopic,
		GroupID: "workflow-group",
	})
	MsgChan = make(chan string)
	Consume()
}

func Consume() {
	go func() {
		for {
			m, err := Reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatal("failed to read:", err)
				return
			}
			MsgChan <- string(m.Value)
		}
	}()
}
