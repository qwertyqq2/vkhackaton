package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

func RsaPublicByte(pub *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pub)
}

func Sign(pk *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, pk, crypto.SHA256, data, nil)
}

func VerifySign(pub *rsa.PublicKey, data, sign []byte) error {
	return rsa.VerifyPSS(pub, crypto.SHA256, data, sign, nil)
}
