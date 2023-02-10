package ring

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

func (sig *Sig) Marshal() ([]byte, error) {
	return json.Marshal(struct {
		R *big.Int
		S *big.Int
	}{
		R: sig.R,
		S: sig.S,
	})
}

func UnmarshalSign(data []byte) (*Sig, error) {
	unmarshalled := struct {
		R *big.Int
		S *big.Int
	}{}
	err := json.Unmarshal(data, &unmarshalled)
	if err != nil {
		return nil, err
	}
	sig := &Sig{}
	sig.R = unmarshalled.R
	sig.S = unmarshalled.S

	return sig, nil
}

func (sig *Sig) Encode() (string, error) {
	b, err := sig.Marshal()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeSign(data string) (*Sig, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	sig, err := UnmarshalSign(b)
	if err != nil {
		return nil, err
	}

	return sig, nil
}
