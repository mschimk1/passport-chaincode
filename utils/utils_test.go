package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {
	msgHash1 := GenerateHash("{\"id\":\"1\", \"from_account\":\"2\", \"to_account\":\"2\", \"amount\":1000, \"fee\":0, \"currency\":\"AUD\"\"}")
	msgHash2 := GenerateHash("{\"id\":\"2\", \"from_account\":\"2\", \"to_account\":\"1\", \"amount\":1000, \"fee\":0, \"currency\":\"AUD\"\"}")
	msgHash3 := GenerateHash("{\"id\":\"1\", \"from_account\":\"2\", \"to_account\":\"2\", \"amount\":1000, \"fee\":0, \"currency\":\"AUD\"\"}")

	assert.Equal(t, msgHash1, msgHash3)
	assert.NotEqual(t, msgHash1, msgHash2)
	assert.NotEqual(t, msgHash2, msgHash3)
}

func TestGenerateID(t *testing.T) {
	assert.Len(t, GenerateID(8), 8)
	assert.NotEqual(t, GenerateID(8), GenerateID(8))
}
