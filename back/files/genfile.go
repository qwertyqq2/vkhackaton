package files

func GenerateFile(data string) *File {
	return &File{
		data: []byte(data),
		Id:   idFile(data),
	}
}
