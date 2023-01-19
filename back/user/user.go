package user

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	mcrypto "github.com/qwertyqq2/filebc/crypto"

	"github.com/qwertyqq2/filebc/values"
)

type User struct {
	pk   *rsa.PrivateKey
	Addr *Address

	Balance uint64
}

func NewUser(pk *rsa.PrivateKey) *User {
	return &User{
		pk:      pk,
		Addr:    NewAddress(pk),
		Balance: uint64(0),
	}
}

func (u *User) Address() *Address {
	return u.Addr
}

func (u *User) Public() string {
	return u.Address().String()
}

func (u *User) SignData(data []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, u.pk, crypto.SHA256, data, nil)
}

func (u *User) Hash() values.Bytes {
	return values.HashSum(u.Addr.Bytes(), mcrypto.ToBytes(u.Balance))
}
