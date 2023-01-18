package crypto

import (
	crand "crypto/rand"
	"crypto/rsa"
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

func GenerateRSAPrivate() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(crand.Reader, 2048)
}
