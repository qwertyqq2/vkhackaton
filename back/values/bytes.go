package values

import (
	"bytes"

	"github.com/qwertyqq2/filebc/crypto"
)

type Bytes []byte

func HashSum(bs ...[]byte) Bytes {
	return crypto.HashSum(
		bytes.Join(
			bs,
			[]byte{},
		))
}
