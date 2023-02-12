package xorstate

import (
	"github.com/qwertyqq2/filebc/values"
)

const (
	LenHash = 32
)

type XorState struct {
	len  int
	seed values.Bytes
}

func NewXorState() *XorState {
	seed := make([]byte, LenHash)
	seed = values.HashSum([]byte("state"))
	return &XorState{
		len:  LenHash,
		seed: seed,
	}
}

func (s *XorState) Add(state values.Bytes, data ...values.Bytes) values.Bytes {
	res := make([]byte, s.len)
	for i, d := range state {
		res[i] = d ^ data[0][i]
	}
	for i := 1; i < len(data); i++ {
		for j, d := range res {
			res[j] = d ^ data[i][j]
		}
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

func (s *XorState) Inverse(data values.Bytes) values.Bytes {
	return data
}
