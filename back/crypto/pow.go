package crypto

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
	"math/rand"
)

func ProowOfWork(blockHash []byte, diff uint8, ch chan bool) ([]byte, bool) {
	var (
		Target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = randuint64()
		hash    []byte
	)
	Target.Lsh(Target, 256-uint(diff))
	for nonce < math.MaxUint64 {
		select {
		case <-ch:
			return nil, false
		default:
			hash = HashSum(bytes.Join(
				[][]byte{
					blockHash,
					ToBytes(nonce),
				},
				[]byte{},
			))
			intHash.SetBytes(hash)
			if intHash.Cmp(Target) == -1 {
				return ToBytes(nonce), true
			}
			nonce += 1
		}
	}
	return nil, true
}

func randuint64() uint64 {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(b)
}
