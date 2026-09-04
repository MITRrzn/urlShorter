package redirect

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
	"urlShorter/internal/helper"
	kafkaProducer "urlShorter/internal/kafka"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{7}$`)

func RedirectHandler(db *sql.DB, redisClient *redis.Client, writer *kafka.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		validationErr := validateCode(code)
		if validationErr != nil {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		valueFromCache, cacheErr := getValueFromCache(r.Context(), redisClient, code)
		if cacheErr == nil {
			handleResolvedLink(w, r, writer, valueFromCache)
			return
		}

		if !errors.Is(cacheErr, redis.Nil) {
			log.Println("redis error:", cacheErr)
		}

		redirectData, err := repository.GetUrlByShortCode(r.Context(), db, code)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helper.WriteErrorResponse(w, "link not found", http.StatusNotFound)
				return
			}

			log.Println(err)
			helper.WriteErrorResponse(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data, marshalErr := json.Marshal(redirectData)
		if marshalErr != nil {
			log.Println(marshalErr)
		} else {
			if cacheSetErr := redisClient.Set(
				r.Context(),
				fmt.Sprintf("link:%s", redirectData.ShortURL),
				data,
				60*time.Minute,
			).Err(); cacheSetErr != nil {
				log.Println("redis set error:", cacheSetErr)
			}
		}

		handleResolvedLink(w, r, writer, redirectData)
	}
}

func handleResolvedLink(w http.ResponseWriter, r *http.Request, writer *kafka.Writer, linkResponse structs.LinkResponse) {
	kafkaProducer.ProcessClickEvent(r, writer, linkResponse)
	http.Redirect(w, r, linkResponse.OriginalURL, http.StatusFound)
}

func getValueFromCache(ctx context.Context, redisClient *redis.Client, key string) (structs.LinkResponse, error) {
	var linkData structs.LinkResponse

	data, err := redisClient.Get(ctx, fmt.Sprintf("link:%s", key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return structs.LinkResponse{}, redis.Nil
	}

	if err != nil {
		return structs.LinkResponse{}, err
	}

	if err := json.Unmarshal(data, &linkData); err != nil {
		return structs.LinkResponse{}, err
	}

	return linkData, nil
}

func validateCode(code string) error {
	if code == "" {
		return errors.New("empty code")
	}
	if len(code) != 7 {
		return errors.New("invalid code")
	}
	if !shortCodeRegex.MatchString(code) {
		return errors.New("incorrect code")
	}

	return nil
}
