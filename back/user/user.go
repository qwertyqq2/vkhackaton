package user

import (
	"crypto/rand"

	mcrypto "github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/crypto/ring"

	"github.com/qwertyqq2/filebc/values"
)

type User struct {
	pk   *ring.PrivateKey
	Addr *Address

	Balance uint64
}

func NewUser(pk *ring.PrivateKey) *User {
	return &User{
		pk:      pk,
		Addr:    NewAddress(pk),
		Balance: uint64(0),
	}
}

func (u *User) Address() *Address {
	return u.Addr
}

func (u *User) Public() *ring.PublicKey {
	return u.Address().Public()
}

func (u *User) Hash() values.Bytes {
	return values.HashSum(u.Addr.Public().Bytes(), mcrypto.ToBytes(u.Balance))
}

func (u *User) RingSignData(data []byte, pubs []*ring.PublicKey, round int) (*ring.Signature, error) {
	sig, err := ring.SignRing(rand.Reader, data, pubs, round, u.pk)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (u *User) SignData(data []byte) (*ring.Sig, error) {
	return ring.SignData(data, u.pk)
}
