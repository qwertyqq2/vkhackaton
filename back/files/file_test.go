package files

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
)

func TestCreateDb(t *testing.T) {
	l, err := NewLevelDB()
	if err != nil {
		t.Fatal(err)
	}
	f1 := NewFile("Здарова епты бля! \n\n\n я знаю")
	rand := []byte("qweqweqweqwe")
	err = l.InsertFile(f1, rand)
	if err != nil {
		t.Fatal(err)
	}
	f2 := NewFile("Че нада")
	err = l.InsertFile(f2, rand)
	if err != nil {
		t.Fatal(err)
	}
	f3 := NewFile("Ниче")
	err = l.InsertFile(f3, rand)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFiles(t *testing.T) {
	l, err := LoadLevel()
	if err != nil {
		t.Fatal(err)
	}
	files, err := l.GetFiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(crypto.Base64EncodeString(f.Id), string(f.Data))
	}
}
