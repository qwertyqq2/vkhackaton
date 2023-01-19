package xorstate

import (
	"bytes"

	"github.com/qwertyqq2/filebc/crypto"
)

type XorState struct {
	len  int
	seed []byte
}

func NewXorState(len int) *XorState {
	seed := make([]byte, len)
	seed = crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte("state"),
			},
			[]byte{},
		))
	return &XorState{
		len:  len,
		seed: seed,
	}
}

func (s *XorState) Add(state []byte, data []byte) []byte {
	res := make([]byte, s.len)
	for i, d := range state {
		res[i] = d ^ data[i]
	}
	return res
}

func (s *XorState) Get(data ...[]byte) []byte {
	res := s.seed
	for _, d := range data {
		res = s.Add(res, d)
	}
	return res
}
