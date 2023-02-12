package xorstate

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/stretchr/testify/assert"
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

	s := NewXorState()

	state := s.Get(id1, id2, id3)
	fmt.Println(crypto.Base64EncodeString(state))
	state = s.Add(state, id4)
	fmt.Println(crypto.Base64EncodeString(state))
	fmt.Println(len(state))
}

func TestInvert(t *testing.T) {
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

	s := NewXorState()

	state1 := s.Get(id1, id2)

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
	invert := s.Inverse(id2)
	if invert == nil {
		t.Fatal(fmt.Errorf("cant invert"))
	}
	state1 = s.Add(state1, invert)
	state1 = s.Add(state1, id3)
	state1 = s.Add(state1, id4)

	state2 := s.Get(id1, id3, id4)

	assert.Equal(t, state1, state2)

}
