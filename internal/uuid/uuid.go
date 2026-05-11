package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateV4 generates a random UUID version 4
func GenerateV4() (string, error) {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", err
	}

	// Set version (4) and variant (RFC 4122)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// GenerateV4Bytes generates a UUID and returns the raw bytes
func GenerateV4Bytes() ([]byte, error) {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return nil, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid[:], nil
}
