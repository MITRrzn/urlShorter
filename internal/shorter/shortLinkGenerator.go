package shorter

import (
	"crypto/rand"
	"math/big"
)

func GenerateShortUrl() (string, error) {
	result := make([]byte, 7)

	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}
