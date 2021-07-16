package hash

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSHA1PasswordHasher_Hash(t *testing.T) {
	h := NewSHA1PasswordHasher("salt")

	password, err := h.Hash("password")

	require.NoError(t, err)
	require.NotNil(t, password)
}
