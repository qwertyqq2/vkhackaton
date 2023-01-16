package files

import (
	"bytes"
	"os"

	"github.com/qwertyqq2/filebc/crypto"
)

const (
	Path = "../filesRepo/"
)

type File struct {
	Id []byte

	FilePath string
	data     string
}

func idFile(data string) string {
	return crypto.Base64EncodeString(crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte(data),
			},
			[]byte{},
		)))
}

func NewFile(data string) (*File, error) {
	id := idFile(data)
	f, err := os.Create(Path + id)
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	return &File{
		data:     data,
		FilePath: Path + id,
		Id:       crypto.Base64DecodeString(id),
	}, nil
}

func verifyId(f *File) bool {
	return bytes.Equal(crypto.Base64DecodeString(string(f.Id)), crypto.Base64DecodeString(idFile(f.data)))
}
