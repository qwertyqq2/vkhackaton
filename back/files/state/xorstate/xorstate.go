package xorstate

import (
	"github.com/qwertyqq2/filebc/values"
)

type XorState struct {
	len  int
	seed values.Bytes
}

func NewXorState(len int) *XorState {
	seed := make([]byte, len)
	seed = values.HashSum([]byte("state"))
	return &XorState{
		len:  len,
		seed: seed,
	}
}

func (s *XorState) Add(state values.Bytes, data values.Bytes) values.Bytes {
	res := make([]byte, s.len)
	for i, d := range state {
		res[i] = d ^ data[i]
	}
	return res
}

func (s *XorState) Get(data ...values.Bytes) values.Bytes {
	res := s.seed
	for _, d := range data {
		res = s.Add(res, d)
	}
	return res
}
