package redirect

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"urlShorter/internal/structs"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValueFromCacheSuccess(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	expected := structs.LinkResponse{
		ID:          1,
		ShortURL:    "QweRty1",
		OriginalURL: "https://example.com",
	}

	data, err := json.Marshal(expected)
	require.NoError(t, err)

	err = mr.Set("link:abc1234", string(data))
	require.NoError(t, err)

	result, err := getValueFromCache(
		context.Background(),
		redisClient,
		"abc1234",
	)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetValueFromCacheNotFound(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	value, getFromCacheErr := getValueFromCache(
		context.Background(),
		redisClient,
		"abc1234",
	)

	assert.ErrorIs(t, getFromCacheErr, redis.Nil)
	assert.Equal(t, structs.LinkResponse{}, value)
}

func TestGetValueFromCacheRedisError(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	mr.Close()

	result, getErr := getValueFromCache(
		context.Background(),
		redisClient,
		"abc1234",
	)

	assert.EqualError(t, getErr, "redis error")
	assert.Equal(t, structs.LinkResponse{}, result)
}

func TestGetValueFromCacheUnmarshalError(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	err = mr.Set("link:abc1234", "{incorrectJSON")
	require.NoError(t, err)

	_, getError := getValueFromCache(
		context.Background(),
		redisClient,
		"abc1234",
	)

	assert.Equal(t, errors.New("unmarshal error"), getError)
}
