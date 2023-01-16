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
	data     []byte
}

func idFile(data string) []byte {
	return crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte(data),
			},
			[]byte{},
		))
}

func NewFile(data string) (*File, error) {
	id := idFile(data)
	path := Path + crypto.Base64EncodeString(id)
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	return &File{
		data:     []byte(data),
		FilePath: path,
		Id:       id,
	}, nil
}

func verifyId(f *File) bool {
	return bytes.Equal(f.Id, idFile(string(f.data)))
}
