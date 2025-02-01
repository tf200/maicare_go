package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63()
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return RandomString(10) + "@example.com"
}

func GetRandomImageURL() string {
	width := 300 // customize size as needed
	height := 200
	id := rand.Intn(1000) // random image ID
	return fmt.Sprintf("https://picsum.photos/id/%d/%d/%d", id, width, height)
}

func RandomPgText() pgtype.Text {
	return pgtype.Text{
		String: RandomString(5),
		Valid:  true,
	}
}

func RandomBool() bool {
	return rand.Float32() < 0.5
}

func RandomPgInt8() pgtype.Int8 {
	return pgtype.Int8{
		Int64: 125,
		Valid: true,
	}
}

func GenerateUsername(firstName, lastName string) string {
	id := uuid.New().String()[:6] // Just 6 characters
	return fmt.Sprintf("%s%s", strings.ToLower(firstName+lastName), id)
}

func RandomTIme() time.Time {
	return time.Now().Add(time.Duration(rand.Int63()))
}
