package files

import (
	"bytes"
	"encoding/json"

	"github.com/qwertyqq2/filebc/crypto"
)

type File struct {
	Id   []byte `json:"id"`
	Data []byte `json:"data"`

	rand []byte
}

func idFile(data string, rand []byte) []byte {
	return crypto.HashSum(
		bytes.Join(
			[][]byte{
				[]byte(data),
				rand,
			},
			[]byte{},
		))
}

func NewFile(data string) *File {
	rand := crypto.GenerateRandom()
	id := idFile(data, rand)

	return &File{
		Data: []byte(data),
		Id:   id,
		rand: rand,
	}
}

func verifyId(f *File) bool {
	return bytes.Equal(f.Id, idFile(string(f.Data), f.rand))
}

func verifySize(f *File, maxSize int) bool {
	if f.size() > maxSize {
		return false
	}
	if f.size() <= 0 {
		return false
	}
	return true
}

func (f *File) Verify(maxSize int) bool {
	if !verifyId(f) || !verifySize(f, maxSize) {
		return false
	}
	return true
}

func (f *File) size() int {
	return len([]rune(string(f.Data)))
}

func (f *File) Diff(maxsize int) int {
	s := f.size()
	return int(s * 100 / maxsize)
}

func (f *File) SerializeFile() (string, error) {
	jsonData, err := json.MarshalIndent(*f, "", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializeFile(fstr string) (*File, error) {
	var f File
	err := json.Unmarshal([]byte(fstr), &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
