package crypto

import (
	"crypto/rsa"
	"crypto/x509"
)

func RsaPublicByte(pub *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pub)
}
