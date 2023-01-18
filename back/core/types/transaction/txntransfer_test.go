package transaction

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/user"

	"github.com/qwertyqq2/filebc/crypto"
)

func TestSerializeTransfer(t *testing.T) {
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u1 := user.NewUser(pk1)
	pk2, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u2 := user.NewUser(pk2)
	tx, err := NewTxTransfer(u1, []byte("first"), u2.Address(), 100)
	if err != nil {
		t.Fatal(err)
	}
	sertx, err := tx.SerializeTx()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sertx)
}

func TestVerifyTransfer(t *testing.T) {
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u1 := user.NewUser(pk1)
	pk2, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u2 := user.NewUser(pk2)
	tx, err := NewTxTransfer(u1, []byte("first"), u2.Address(), 100)
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
