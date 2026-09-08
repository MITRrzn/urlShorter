package shorter

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"urlShorter/internal/structs"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateLinkHandlerSuccess(t *testing.T) {
	oldStore := storeShortURLFn
	t.Cleanup(func() {
		storeShortURLFn = oldStore
	})

	testCases := []struct {
		name           string
		body           string
		storeResult    string
		storeErr       error
		expectedStatus int
	}{
		{
			name:           "success",
			body:           `{"url":"https://example.com"}`,
			storeResult:    "abc1234",
			storeErr:       nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{"url":`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid url",
			body:           `{"url":"example.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "store error",
			body:           `{"url":"https://example.com"}`,
			storeErr:       errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storeShortURLFn = func(ctx context.Context, db *sql.DB, link structs.LinkStruct) (string, error) {
				return testCase.storeResult, testCase.storeErr
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/links",
				strings.NewReader(testCase.body),
			)

			recorder := httptest.NewRecorder()
			handler := CreateLinkHandler(nil)
			handler.ServeHTTP(recorder, req)
			if testCase.expectedStatus != recorder.Code {
				t.Errorf("expected status %d, got %d", testCase.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestStoreShortURLRetryOnDuplicate(t *testing.T) {
	oldGenerate := generateShortURL
	oldAddLink := addLink
	t.Cleanup(func() {
		generateShortURL = oldGenerate
		addLink = oldAddLink
	})

	generateCounter := 0
	addLinkCounter := 0

	generateShortURL = func() (string, error) {
		generateCounter++

		if generateCounter == 1 {
			return "first12", nil
		}

		return "second1", nil
	}

	addLink = func(ctx context.Context, db *sql.DB, originalURL string, shortURL string) error {
		addLinkCounter++
		if addLinkCounter == 1 {
			return &pq.Error{
				Code: "23505",
			}
		}

		return nil
	}

	link := structs.LinkStruct{
		URL: "https://example.com",
	}

	result, err := storeShortURL(context.Background(), nil, link)

	assert.NoError(t, err)
	assert.Equal(t, "second1", result)
	assert.Equal(t, 2, generateCounter)
	assert.Equal(t, 2, addLinkCounter)
}

func TestStoreShortURLGenerateError(t *testing.T) {
	oldGenerate := generateShortURL
	t.Cleanup(func() {
		generateShortURL = oldGenerate
	})

	expectedErr := errors.New("generation failed")

	generateShortURL = func() (string, error) {
		return "", expectedErr
	}

	link := structs.LinkStruct{
		URL: "https://example.com",
	}

	result, err := storeShortURL(context.Background(), nil, link)

	assert.Empty(t, result)
	assert.Equal(t, expectedErr, err)
}

func TestStoreShortURLSuccess(t *testing.T) {
	oldGenerate := generateShortURL
	oldAddLink := addLink
	t.Cleanup(func() {
		generateShortURL = oldGenerate
		addLink = oldAddLink
	})

	generateShortURL = func() (string, error) {
		return "qwer123", nil
	}

	addLink = func(ctx context.Context, db *sql.DB, originalURL string, shortURL string) error {
		return nil
	}

	link := structs.LinkStruct{
		URL: "https://example.com",
	}

	result, err := storeShortURL(context.Background(), nil, link)

	assert.NoError(t, err)
	assert.Equal(t, "qwer123", result)
}

func TestValidateURL(t *testing.T) {
	testCases := []struct {
		name     string
		URL      string
		expected error
	}{
		{
			name:     "Validation correct url",
			URL:      "https://example.com/my/path?param=value",
			expected: nil,
		},
		{
			name:     "validation incorrect url, missed scheme",
			URL:      "example.com/my/path?param=value",
			expected: errors.New("scheme is empty"),
		},
		{
			name:     "validation incorrect url, missed host",
			URL:      "https://",
			expected: errors.New("host is empty"),
		},
		{
			name:     "validation incorrect url, incorrect schema",
			URL:      "mailto://example.com/my/path?param=value",
			expected: errors.New("url scheme must be http or https"),
		},
		{
			name:     "validation empty url",
			URL:      "",
			expected: errors.New("url is empty"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := validateURL(testCase.URL)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
