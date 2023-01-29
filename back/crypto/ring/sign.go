package ring

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

type Sig struct {
	R *big.Int
	S *big.Int
}

func SignData(data []byte, pk *PrivateKey) (*Sig, error) {
	pkecdsa := pk.ToEcdsa()
	if pkecdsa == nil {
		return nil, fmt.Errorf("cant get ecdsa")
	}
	r, s, err := ecdsa.Sign(rand.Reader, pkecdsa, data)
	if err != nil {
		return nil, err
	}
	return &Sig{
		R: r,
		S: s,
	}, nil
}

func VerifySign(data []byte, sig *Sig, pub *PublicKey) bool {
	pubecdsa := pub.ToEcdsa()
	if pubecdsa == nil {
		return false
	}
	if sig.S.Cmp(big.NewInt(0)) == -1 || sig.R.Cmp(big.NewInt(0)) == -1 {
		return false
	}
	return ecdsa.Verify(pubecdsa, data, sig.R, sig.S)
}
