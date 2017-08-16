package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	assert.Len(t, GenerateID(8), 8)
	assert.NotEqual(t, GenerateID(8), GenerateID(8))
}
