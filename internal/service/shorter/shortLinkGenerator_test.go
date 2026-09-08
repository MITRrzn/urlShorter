package shorter

import (
	"fmt"
	"testing"
	"urlShorter/internal/helper"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortUrl(t *testing.T) {
	shortURL, err := GenerateShortUrl()
	if err != nil {
		t.Error(err)
	}

	isValid := helper.ValidateCode(shortURL)
	if !assert.Nil(t, isValid) {
		fmt.Println(shortURL, isValid)
		t.Error("generated shortURL is not valid")
	}
}
