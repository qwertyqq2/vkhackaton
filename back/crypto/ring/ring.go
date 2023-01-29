package ring

import (
	"bytes"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math/big"
)

var (
	ErrEmptyMessage = errors.New("you should provide a message to sign")

	ErrInvalidSignerIndex = errors.New("the index of the signer should be in the ring")

	ErrRingTooSmall = errors.New("the ring is too small: you need at least two participants")
)

type Signature struct {
	Ring  []*PublicKey
	Seed  []byte
	Sings [][]byte
}

func SignRing(
	rand io.Reader,
	message []byte,
	ringKeys []*PublicKey,
	round int,
	sk *PrivateKey,
) (*Signature, error) {
	if len(message) == 0 {
		return nil, ErrEmptyMessage
	}
	if round < 0 || len(ringKeys) <= round {
		return nil, ErrInvalidSignerIndex
	}
	if len(ringKeys) < 2 {
		return nil, ErrRingTooSmall
	}
	if rand == nil {
		rand = crand.Reader
	}
	es := make([][]byte, len(ringKeys))
	ss := make([][]byte, len(ringKeys))
	curve := elliptic.P384()
	r := len(ringKeys)
	k, err := randomParam(curve, rand)
	if err != nil {
		return nil, err
	}
	x, y := curve.ScalarBaseMult(k)
	es[(round+1)%r] = hash(append(message, elliptic.Marshal(curve, x, y)...))
	for i := (round + 1) % r; i != round; i = (i + 1) % r {
		s, err := randomParam(curve, rand)
		if err != nil {
			return nil, err
		}
		ss[i] = s
		x1, y1 := curve.ScalarBaseMult(ss[i])
		px, py := ringKeys[i].x, ringKeys[i].y
		x2, y2 := curve.ScalarMult(px, py, es[i])
		x, y = curve.Add(x1, y1, x2, y2)
		es[(i+1)%r] = hash(append(message, elliptic.Marshal(curve, x, y)...))
	}
	valK := new(big.Int).SetBytes(k)
	valE := new(big.Int).SetBytes(es[round])
	valX := new(big.Int).SetBytes(sk.d)
	valS := new(big.Int).Sub(valK, new(big.Int).Mul(valE, valX))
	if valS.Sign() == -1 {
		add := new(big.Int).Mul(valE, curve.Params().N)
		valS = valS.Add(valS, add)
		_, valS = new(big.Int).DivMod(valS, curve.Params().N, new(big.Int))
		if valS.Sign() == 0 {
			return nil, errors.New("could not produce ring signature")
		}
	}
	ss[round] = valS.Bytes()
	sig := &Signature{
		Ring:  ringKeys,
		Seed:  es[0],
		Sings: ss,
	}
	return sig, nil
}

func randomParam(curve elliptic.Curve, rand io.Reader) ([]byte, error) {
	for {
		r, err := crand.Int(rand, curve.Params().N)
		if err != nil {
			return nil, err
		}
		if r.Sign() == 1 {
			return r.Bytes(), nil
		}
	}
}

func hash(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func VerifyRing(message []byte, sig *Signature) bool {
	if sig == nil {
		return false
	}
	if len(sig.Ring) < 2 {
		return false
	}
	if len(sig.Sings) != len(sig.Ring) {
		return false
	}
	if len(sig.Seed) == 0 {
		return false
	}
	curve := elliptic.P384()
	e := make([]byte, len(sig.Seed))
	copy(e, sig.Seed)
	for i := 0; i < len(sig.Ring); i++ {
		x1, y1 := curve.ScalarBaseMult(sig.Sings[i])
		px, py := sig.Ring[i].x, sig.Ring[i].y
		x2, y2 := curve.ScalarMult(px, py, e)
		x, y := curve.Add(x1, y1, x2, y2)
		e = hash(append(message, elliptic.Marshal(curve, x, y)...))
	}
	return bytes.Equal(e, sig.Seed)
}
