package token_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zerok/remarked/internal/token"
)

func TestGenerateToken(t *testing.T) {
	// token.Generate should generate a random string with 6 characters
	tkn := token.Generate()
	assert.Len(t, tkn, 6, "The generated token should have 6 characters")
}
