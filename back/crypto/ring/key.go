package ring

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

type PrivateKey struct {
	d []byte
}

type PublicKey struct {
	x     *big.Int
	y     *big.Int
	curve elliptic.Curve
}

func (pk *PrivateKey) Public() *PublicKey {
	var priv ecdsa.PrivateKey
	priv.D = new(big.Int).SetBytes(pk.d)
	curve := elliptic.P384()
	priv.PublicKey.Curve = curve
	X, Y := priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())
	return &PublicKey{
		curve: curve,
		x:     X,
		y:     Y,
	}
}

func (pk *PrivateKey) String() string {
	return base64.StdEncoding.EncodeToString(pk.d)
}

func (pub *PublicKey) Bytes() []byte {
	return elliptic.Marshal(pub.curve, pub.x, pub.y)

}

func (pub *PublicKey) String() string {
	return base64.StdEncoding.EncodeToString(pub.Bytes())
}

func ParsePublic(pubs string) *PublicKey {
	pubB, err := base64.StdEncoding.DecodeString(pubs)
	if err != nil {
		return nil
	}
	x, y := elliptic.Unmarshal(elliptic.P384(), pubB)
	return &PublicKey{
		curve: elliptic.P384(),
		x:     x,
		y:     y,
	}
}

func ParsePrivate(pks string) *PrivateKey {
	d, err := base64.StdEncoding.DecodeString(pks)
	if err != nil {
		return nil
	}
	return &PrivateKey{
		d: d,
	}
}

func (pk *PrivateKey) Marshal() []byte {
	return pk.d
}

func UnmarshalPrivate(d []byte) *PrivateKey {
	return &PrivateKey{
		d: d,
	}
}

func GeneratePrivate() *PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil
	}
	return &PrivateKey{
		d: privateKey.D.Bytes(),
	}
}

func (pub *PublicKey) ToEcdsa() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		X:     pub.x,
		Y:     pub.y,
		Curve: pub.curve,
	}
}

func (pk *PrivateKey) ToEcdsa() *ecdsa.PrivateKey {
	pub := pk.Public()
	if pub == nil {
		return nil
	}
	return &ecdsa.PrivateKey{
		D:         new(big.Int).SetBytes(pk.d),
		PublicKey: *pub.ToEcdsa(),
	}
}
