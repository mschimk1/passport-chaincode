package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

// GenerateHash creates a sha256 hash of the given message
func GenerateHash(msg string) string {
	msgHash := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", msgHash)
}

// GenerateID generates a fixed length string of random digits
func GenerateID(length int) string {
	rand.Seed(time.Now().UnixNano())
	r := []rune("1234567890")
	b := make([]rune, length)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}
