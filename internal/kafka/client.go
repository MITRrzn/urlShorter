package kafka

import (
	"log"

	"github.com/segmentio/kafka-go"
)

func GetKafkaClient() *kafka.Writer {
	writer := &kafka.Writer{
		Addr:  kafka.TCP("kafka:9092"),
		Topic: "urlShorter-clicks",
		Async: true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Printf("kafka write error: %v", err)
			}
		},
	}
	defer func(writer *kafka.Writer) {
		writerErr := writer.Close()
		if writerErr != nil {
			log.Println(writerErr)
		}
	}(writer)

	return writer
}
