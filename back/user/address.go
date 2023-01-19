package user

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"

	mcrypto "github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/values"
)

type Address struct {
	pub *rsa.PublicKey
}

func NewAddress(pk *rsa.PrivateKey) *Address {
	return &Address{
		pub: &pk.PublicKey,
	}
}

func (a *Address) String() string {
	return mcrypto.Base64EncodeString(x509.MarshalPKCS1PublicKey(a.pub))
}

func (a *Address) ToBytes() ([]byte, error) {
	spub, err := x509.MarshalPKIXPublicKey(a.pub)
	if err != nil {
		return nil, err
	}
	return spub, nil
}

func VerifySign(a *Address, data, sign []byte) error {
	return rsa.VerifyPSS(a.pub, crypto.SHA256, data, sign, nil)
}

func ParseAddress(saddr string) *Address {
	addrb := mcrypto.Base64DecodeString(saddr)
	pub, err := x509.ParsePKCS1PublicKey(addrb)
	if err != nil {
		return nil
	}
	return &Address{
		pub: pub,
	}
}

func (a *Address) Bytes() values.Bytes {
	return mcrypto.RsaPublicByte(a.pub)
}
