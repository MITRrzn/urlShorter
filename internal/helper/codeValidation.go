package helper

import (
	"errors"
	"regexp"
)

var shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{7}$`)

func ValidateCode(code string) error {
	if !shortCodeRegex.MatchString(code) {
		return errors.New("incorrect code")
	}
	return nil
}
