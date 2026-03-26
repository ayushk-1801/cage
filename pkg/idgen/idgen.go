package idgen

import (
	"crypto/rand"
	"encoding/hex"
)

func Generate() string {
	b := make([]byte, 6) // 12 hex chars
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}