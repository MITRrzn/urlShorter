package redirect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"

	"github.com/redis/go-redis/v9"
)

func RedirectHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		validationErr := validateCode(code)
		if validationErr != nil {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		valueFromCache, cacheErr := getValueFromCache(r.Context(), redisClient, code)
		if cacheErr == nil {
			log.Println("get value from cache", valueFromCache)
			http.Redirect(w, r, valueFromCache, http.StatusFound)
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

		if cacheSetErr := redisClient.Set(
			r.Context(),
			fmt.Sprintf("link:%s", redirectData.ShortURL),
			redirectData.OriginalURL,
			60*time.Minute,
		).Err(); cacheSetErr != nil {
			log.Println("redis set error:", cacheSetErr)
		}

		http.Redirect(w, r, redirectData.OriginalURL, http.StatusFound)
	}
}

func getValueFromCache(ctx context.Context, redisClient *redis.Client, key string) (string, error) {
	val, getErr := redisClient.Get(ctx, fmt.Sprintf("link:%s", key)).Result()

	if errors.Is(getErr, redis.Nil) {
		return "", redis.Nil
	}

	if getErr != nil {
		fmt.Printf("failed to get value, error: %v\n", getErr)
		return "", getErr
	}

	return val, nil
}

func validateCode(code string) error {
	if code == "" {
		return errors.New("empty code")
	}
	if len(code) != 7 {
		return errors.New("invalid code")
	}
	var shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{7}$`)
	if !shortCodeRegex.MatchString(code) {
		return errors.New("incorrect code")
	}

	return nil
}
