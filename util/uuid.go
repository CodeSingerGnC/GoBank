package util

import ( 
	"github.com/google/uuid"
)

func BytesToUUID(b []byte) (uuid.UUID, error) {
	return uuid.FromBytes(b)
}

func UUIDToBytes(u uuid.UUID) []byte {
	return u[:]
}