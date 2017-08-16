package utils

import (
	"math/rand"
	"time"
)

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
