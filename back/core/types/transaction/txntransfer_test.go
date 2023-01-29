package transaction

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/user"
)

func TestSerializeTransfer(t *testing.T) {
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
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
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
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
