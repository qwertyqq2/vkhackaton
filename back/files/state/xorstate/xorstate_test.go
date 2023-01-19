package xorstate

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
)

func TestState(t *testing.T) {
	id1 := crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte("qweqweqeqew"),
			},
			[]byte{},
		))
	id2 := crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte("asdazxvxzv"),
			},
			[]byte{},
		))
	id3 := crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte("dgdfgdfgdfas"),
			},
			[]byte{},
		))
	id4 := crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte("zvxcvcx"),
			},
			[]byte{},
		))

	s := NewXorState(len(id1))

	state := s.Get(id1, id2, id3)
	fmt.Println(crypto.Base64EncodeString(state))
	state = s.Add(state, id4)
	fmt.Println(crypto.Base64EncodeString(state))
	fmt.Println(len(state))
}
