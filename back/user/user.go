package user

import (
	"crypto/rand"
	"encoding/json"

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

func GetUser(a *Address, bal uint64) *User {
	return &User{
		Addr:    a,
		Balance: bal,
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

func (u *User) Id() values.Bytes {
	return u.Hash()
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

func (u *User) Serialize() (string, error) {
	jsonData, err := json.MarshalIndent(u, " ", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func Deserialize(data string) (*User, error) {
	var u User
	err := json.Unmarshal([]byte(data), &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
