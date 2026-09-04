package kafka

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"urlShorter/internal/structs"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func ProcessClickEvent(r *http.Request, writer *kafka.Writer, linkResponse structs.LinkResponse) {
	clickData := structs.ClickEvent{
		EventID:   uuid.New().String(),
		LinkID:    linkResponse.ID,
		ShortCode: linkResponse.ShortURL,
		ClickedAt: time.Now().UTC(),
		Referer:   r.Referer(),
		UserAgent: r.UserAgent(),
	}

	data, err := json.Marshal(clickData)
	if err != nil {
		log.Println("json marshal error:", err)
		return
	}

	err = writer.WriteMessages(r.Context(), kafka.Message{
		Key:   []byte(linkResponse.ShortURL),
		Value: data,
	})
	if err != nil {
		log.Println("write kafka messages error:", err)
		return
	}
}
