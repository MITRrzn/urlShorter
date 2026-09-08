package shorter

import (
	"testing"
	"urlShorter/internal/helper"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortUrl(t *testing.T) {
	shortURL, err := GenerateShortUrl()
	if err != nil {
		t.Fatal(err)
	}

	isValid := helper.ValidateCode(shortURL)
	assert.Nil(t, isValid)
}
