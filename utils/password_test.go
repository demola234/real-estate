package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(8)

	hashed, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashed)
}

func TestVerifyPassword(t *testing.T) {
	password := RandomString(8)

	hashed, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	check, msg := VerifyPassword(password, hashed)

	require.True(t, check)
	require.Empty(t, msg)
}
