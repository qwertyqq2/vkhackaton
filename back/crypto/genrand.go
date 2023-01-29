package crypto

import (
	"math/rand"
	"time"
)

func GenerateRandom() []byte {
	rand.Seed(time.Now().UnixNano())
	slice := make([]byte, 100)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}
