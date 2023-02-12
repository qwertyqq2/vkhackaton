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

// func Marhal(msg *Message) ([]byte, error) {
// 	return json.Marshal(struct {
// 		Id      int
// 		Payload []byte
// 	}{
// 		Id:      msg.Id,
// 		Payload: msg.payload,
// 	})
// }

// func Unmarhsal(d []byte) (*Message, error) {
// 	unmarshalled := struct {
// 		Id      int
// 		Payload []byte
// 	}{}
// 	err := json.Unmarshal(d, &unmarshalled)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Message{
// 		Id:      unmarshalled.Id,
// 		payload: unmarshalled.Payload,
// 	}, nil
// }
