package user

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

type User struct {
	pk   *rsa.PrivateKey
	addr Address
}

func NewUser(pk *rsa.PrivateKey) *User {
	return &User{
		pk:   pk,
		addr: *NewAddress(pk),
	}
}

func (u *User) Address() *Address {
	return &u.addr
}

func (u *User) Public() string {
	return u.Address().String()
}

func (u *User) SignData(data []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, u.pk, crypto.SHA256, data, nil)
}
