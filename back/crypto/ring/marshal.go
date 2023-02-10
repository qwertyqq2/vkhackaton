package ring

import (
	"encoding/base64"
	"encoding/json"
)

func (sig *Signature) Marshal() ([]byte, error) {
	return json.Marshal(struct {
		R []*PublicKey
		S [][]byte
		E []byte
	}{
		R: sig.Ring,
		S: sig.Sings,
		E: sig.Seed,
	})
}

func UnmarshalRing(data []byte) (*Signature, error) {
	unmarshalled := struct {
		R []*PublicKey
		S [][]byte
		E []byte
	}{}
	err := json.Unmarshal(data, &unmarshalled)
	if err != nil {
		return nil, err
	}
	sig := &Signature{}
	sig.Ring = unmarshalled.R
	sig.Seed = unmarshalled.E
	sig.Sings = unmarshalled.S

	return sig, nil
}

func (sig *Signature) Encode() (string, error) {
	b, err := sig.Marshal()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeRing(data string) (*Signature, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	sig, err := UnmarshalRing(b)
	if err != nil {
		return nil, err
	}

	return sig, nil
}
