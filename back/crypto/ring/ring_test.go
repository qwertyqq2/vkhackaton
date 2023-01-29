package ring

import (
	"fmt"
	"log"
	"testing"
)

func TestRing(t *testing.T) {
	var N = 10
	pubKeys := make([]*PublicKey, 0)
	var secr *PrivateKey
	idx := 5
	for i := 0; i < N; i++ {
		pk := GeneratePrivate()
		pub := pk.Public()
		pubKeys = append(pubKeys, pub)
		if i == idx {
			secr = pk
		}
	}
	sign, err := SignRing(nil, []byte("aaaaa"), pubKeys, idx, secr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(VerifyRing([]byte("aaaaa"), sign))
}
