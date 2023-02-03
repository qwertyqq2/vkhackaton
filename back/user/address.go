package user

import (
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/qwertyqq2/filebc/crypto/ring"
)

type Address struct {
	pub *ring.PublicKey
}

func NewAddress(pk *ring.PrivateKey) *Address {
	return &Address{
		pub: pk.Public(),
	}
}

func (a *Address) String() string {
	return a.pub.String()
}

func (a *Address) Public() *ring.PublicKey {
	return a.pub
}

func VeryfySignRing(data []byte, addr []*Address, seed []byte, signs [][]byte) bool {
	pubs := make([]*ring.PublicKey, len(addr))
	for i := 0; i < len(addr); i++ {
		pubs[i] = addr[i].Public()
	}
	return ring.VerifyRing(data, &ring.Signature{
		Ring:  pubs,
		Seed:  seed,
		Sings: signs,
	})
}

func ParseAddress(saddr string) (*Address, error) {
	pub := ring.ParsePublic(saddr)
	if pub == nil {
		return nil, fmt.Errorf("nil pub")
	}
	return &Address{
		pub: pub,
	}, nil
}

func Shuffle(a []string) {
	mrand.Seed(time.Now().UnixNano())
	mrand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}
