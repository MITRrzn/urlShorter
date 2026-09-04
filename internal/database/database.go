package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func PsqlConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if pingErr := db.Ping(); pingErr != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, pingErr
	}

	return db, nil
}

func GetRedisClient(ctx context.Context) (*redis.Client, error) {
	redisDB, envErr := strconv.Atoi(os.Getenv("REDIS_DB"))
	if envErr != nil {
		log.Println(envErr)
		return nil, envErr
	}

	db := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"),
		Username:     os.Getenv("REDIS_USER"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           redisDB,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		MaxRetries:   2,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
		closeErr := db.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	log.Printf("success redis connect")

	return db, nil
}
