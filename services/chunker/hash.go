package chunker

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashChunk calculates SHA-256 hash of given bytes and returns it as hex string
func HashChunk(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
