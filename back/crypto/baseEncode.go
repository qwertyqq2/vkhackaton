package crypto

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

func Base64EncodeString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64DecodeString(data string) []byte {
	result, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return result
}

func ToBytes(num uint64) []byte {
	data := new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}
	return data.Bytes()
}
