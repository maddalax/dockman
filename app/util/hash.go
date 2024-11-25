package util

import (
	"crypto/sha256"
	"fmt"
)

func HashString(data string) string {
	// Compute the SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hex string
	return fmt.Sprintf("%x", hash)
}
