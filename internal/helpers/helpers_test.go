package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortCode(t *testing.T) {
	originalURL := "https://test.test"
	length := 6

	shortCode := GenerateShortCode(originalURL, length)
	assert.Equal(t, length, len(shortCode))
	assert.Equal(t, GenerateShortCode(originalURL, length), shortCode, "Short code generation should be deterministic")
	assert.NotEqual(t, length, GenerateShortCode(originalURL, 10), "Short code length should not be equal to the length parameter")
}
