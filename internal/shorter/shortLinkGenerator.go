package shorter

import (
	"crypto/rand"
	"math/big"
)

func GenerateShortUrl() (string, error) {
	result := make([]byte, 7)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}
