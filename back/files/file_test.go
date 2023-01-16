package files

import "testing"

func TestCreateFile(t *testing.T) {
	data := "hello\ngo\n"

	_, err := NewFile(data)
	if err != nil {
		t.Log(err)
	}

}
