package values

import (
	"bytes"
	"crypto/sha256"
)

type Bytes []byte

func HashSum(bs ...[]byte) Bytes {
	return hashSum(
		bytes.Join(
			bs,
			[]byte{},
		))
}

func hashSum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
