package crypto

import "encoding/base64"

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
