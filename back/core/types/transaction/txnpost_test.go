package transaction

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
)

func TestSerializePost(t *testing.T) {
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u1 := user.NewUser(pk1)
	file := files.GenerateFile("first fileqweqweqwweqwwwwwwwwwwwwwwwwwwwwwwwwsqdqsdqwdqdwqddwq")
	tx, err := NewTxPost(u1, []byte("first"), file)
	if err != nil {
		t.Fatal(err)
	}
	sertx, err := tx.SerializeTx()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sertx)
}

func TestVerifyPost(t *testing.T) {
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u1 := user.NewUser(pk1)
	file := files.GenerateFile("first fileqweqweqwweqwwwwwwwwwwwwwwwwwwwwwwwwsqdqsdqwdqdwqddwq")
	tx, err := NewTxPost(u1, []byte("first"), file)
	if err != nil {
		t.Fatal(err)
	}
	f := tx.Valid()
	if !f {
		t.Log("NOT VERIFY")
	} else {
		t.Log("VERIFY")
	}
}
