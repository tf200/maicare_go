package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Helper function to generate random recovery codes
func GenerateRecoveryCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		// Generate a code like "XXXX-XXXX-XXXX" (hyphenated for readability)
		code := fmt.Sprintf("%s-%s-%s",
			randomCode(4),
			randomCode(4),
			randomCode(4))
		codes[i] = code
	}
	return codes
}

func randomCode(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed similar-looking characters
	b := make([]byte, length)
	for i := range b {
		// Generate a secure random index
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}
