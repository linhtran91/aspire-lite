package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	secretKey := "9876hjds"
	duration := 5 * time.Minute
	builder := NewJWTTokenBuilder(secretKey, duration)

	token, err := builder.Encode(int64(1))
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestDecode(t *testing.T) {
	secretKey := "9876hjds"
	duration := 5 * time.Minute
	builder := NewJWTTokenBuilder(secretKey, duration)

	customerID := int64(23)
	token, err := builder.Encode(customerID)
	assert.NoError(t, err)

	got, err := builder.Decode(token)
	assert.NoError(t, err)
	assert.Equal(t, customerID, got)
}
