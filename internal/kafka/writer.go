package kafka

import (
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func NewWriter() *kafka.Writer {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic: "urlShorter-clicks",
		Async: true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Printf("kafka write error: %v", err)
			}
		},
	}

	return writer
}
