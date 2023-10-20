package id

import (
	"math/rand"
	"time"

	"github.com/jakehl/goid"
)

// GetUUIDV1 ...
func GetUUIDV1() (uid string) {
	var u *goid.UUID
	u = goid.NewV4UUID()
	uid = u.String()
	return
}

// GetUUIDV2 ...
func GetUUIDV2(length int) (uid string) {
	rand.Seed(time.Now().UTC().UnixNano())
	time.Sleep(time.Nanosecond)

	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	uid = string(b)
	return
}
