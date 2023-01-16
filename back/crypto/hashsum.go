package crypto

import "crypto/sha256"

func HashSum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
