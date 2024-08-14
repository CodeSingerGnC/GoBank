package util

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUUIDToBytes(t *testing.T) {
	u := uuid.New()

	b := UUIDToBytes(u)
	require.Equal(t, 16, len(b))

	u2, err := BytesToUUID(b)
	require.NoError(t, err)
	require.Equal(t, u, u2)
}

func TestBytesToUUID(t *testing.T) {
	u := uuid.New()
	b := UUIDToBytes(u)

	u2, err := BytesToUUID(b)
	require.NoError(t, err)
	require.Equal(t, u, u2)
}