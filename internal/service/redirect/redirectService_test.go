package redirect

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlShorter/internal/structs"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
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

	assert.EqualError(t, getError, "unmarshal error")
}

func TestSetValueToCacheSuccess(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	data := structs.LinkResponse{
		ID:          1,
		ShortURL:    "QweRty1",
		OriginalURL: "https://example.com",
	}

	err = setValueToCacheFn(
		context.Background(),
		redisClient,
		data,
	)
	require.NoError(t, err)

	cachedValue, err := mr.Get("link:QweRty1")
	require.NoError(t, err)
	valueFromCache := json.Unmarshal([]byte(cachedValue), &structs.LinkResponse{})
	assert.Equal(t, data, valueFromCache)
}

func TestSetValueToCacheSetErr(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	mr.Close()

	data := structs.LinkResponse{
		ID:          1,
		ShortURL:    "QweRty1",
		OriginalURL: "https://example.com",
	}

	err = setValueToCacheFn(
		context.Background(),
		redisClient,
		data,
	)

	assert.Error(t, err)
}

func TestRedirectHandlerCacheHit(t *testing.T) {
	oldGetCache := getValueFromCacheFn
	oldHandle := handleResolvedLinkFn

	t.Cleanup(func() {
		getValueFromCacheFn = oldGetCache
		handleResolvedLinkFn = oldHandle
	})

	expected := structs.LinkResponse{
		ID:          1,
		ShortURL:    "QweRty1",
		OriginalURL: "https://example.com",
	}

	getValueFromCacheFn = func(
		ctx context.Context,
		redisClient *redis.Client,
		key string,
	) (structs.LinkResponse, error) {
		return expected, nil
	}

	handleCalls := 0

	handleResolvedLinkFn = func(
		w http.ResponseWriter,
		r *http.Request,
		writer *kafka.Writer,
		data structs.LinkResponse,
	) {
		handleCalls++
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/QweRty1",
		nil,
	)

	req.SetPathValue("code", "QweRty1")

	recorder := httptest.NewRecorder()

	handler := RedirectHandler(nil, nil, nil)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, 1, handleCalls)
}

func TestRedirectHandlerDbHit(t *testing.T) {
	oldHandle := handleResolvedLinkFn
	oldGetCache := getValueFromCacheFn
	oldGetURLByShortCode := getURLByShortCodeFn
	oldSetValueToCacheFn := setValueToCacheFn

	t.Cleanup(func() {
		handleResolvedLinkFn = oldHandle
		getURLByShortCodeFn = oldGetURLByShortCode
		getValueFromCacheFn = oldGetCache
		setValueToCacheFn = oldSetValueToCacheFn
	})

	expected := structs.LinkResponse{
		ID:          1,
		ShortURL:    "QweRty1",
		OriginalURL: "https://example.com",
	}

	getValueFromCacheFn = func(
		ctx context.Context,
		redisClient *redis.Client,
		key string,
	) (structs.LinkResponse, error) {
		return structs.LinkResponse{}, redis.Nil
	}

	dbCalls := 0
	getURLByShortCodeFn = func(
		ctx context.Context,
		db *sql.DB,
		code string,
	) (structs.LinkResponse, error) {
		dbCalls++
		return expected, nil
	}

	handleCalls := 0
	var actual structs.LinkResponse
	handleResolvedLinkFn = func(
		w http.ResponseWriter,
		r *http.Request,
		writer *kafka.Writer,
		data structs.LinkResponse,
	) {
		actual = data
		handleCalls++
	}

	cacheSetCalls := 0
	setValueToCacheFn = func(
		ctx context.Context,
		redisClient *redis.Client,
		redirectData structs.LinkResponse,
	) error {
		cacheSetCalls++
		return nil
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/QweRty1",
		nil,
	)

	req.SetPathValue("code", "QweRty1")

	recorder := httptest.NewRecorder()

	handler := RedirectHandler(nil, nil, nil)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, 1, dbCalls)
	assert.Equal(t, 1, cacheSetCalls)
	assert.Equal(t, 1, handleCalls)
	assert.Equal(t, expected, actual)
}
