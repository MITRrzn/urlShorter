package redirect

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"urlShorter/internal/helper"
	kafkaProducer "urlShorter/internal/kafka"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var getValueFromCacheFn = getValueFromCache
var getURLByShortCodeFn = repository.GetUrlByShortCode
var handleResolvedLinkFn = handleResolvedLink
var setValueToCacheFn = setValueToCache
var validateCodeFn = helper.ValidateCode

func RedirectHandler(db *sql.DB, redisClient *redis.Client, writer *kafka.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		validationErr := validateCodeFn(code)
		if validationErr != nil {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		valueFromCache, cacheErr := getValueFromCacheFn(r.Context(), redisClient, code)
		if cacheErr == nil {
			handleResolvedLinkFn(w, r, writer, valueFromCache)
			return
		}

		if !errors.Is(cacheErr, redis.Nil) {
			log.Println("redis error:", cacheErr)
		}

		redirectData, err := getURLByShortCodeFn(r.Context(), db, code)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helper.WriteErrorResponse(w, "link not found", http.StatusNotFound)
				return
			}

			log.Println(err)
			helper.WriteErrorResponse(w, "internal server error", http.StatusInternalServerError)
			return
		}

		setErr := setValueToCacheFn(r.Context(), redisClient, redirectData)
		if setErr != nil {
			log.Println("redis set error:", setErr)
		}

		handleResolvedLinkFn(w, r, writer, redirectData)
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
		return structs.LinkResponse{}, errors.New("redis error")
	}

	if unmarshalErr := json.Unmarshal(data, &linkData); unmarshalErr != nil {
		return structs.LinkResponse{}, errors.New("unmarshal error")
	}

	return linkData, nil
}

func setValueToCache(ctx context.Context, redisClient *redis.Client, redirectData structs.LinkResponse) error {
	data, err := json.Marshal(redirectData)
	if err != nil {
		return err
	}

	return redisClient.Set(
		ctx,
		fmt.Sprintf("link:%s", redirectData.ShortURL),
		data,
		60*time.Minute,
	).Err()
}
