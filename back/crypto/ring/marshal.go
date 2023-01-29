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

func (sig *Signature) Unmarshal(data []byte) error {
	unmarshalled := struct {
		R []*PublicKey
		S [][]byte
		E []byte
	}{}
	err := json.Unmarshal(data, &unmarshalled)
	if err != nil {
		return err
	}

	sig.Ring = unmarshalled.R
	sig.Seed = unmarshalled.E
	sig.Sings = unmarshalled.S

	return nil
}

func (sig *Signature) Encode() (string, error) {
	b, err := sig.Marshal()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (sig *Signature) Decode(data string) error {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	err = sig.Unmarshal(b)
	if err != nil {
		return err
	}

	return nil
}
