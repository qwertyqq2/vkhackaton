package types

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
)

func CreatePostTx() (string, error) {
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	singers := []*user.Address{u2.Addr, u3.Addr}

	file := files.NewFile("first fileqweqweqwweqwwwwwwwwwwwwwwwwwwwwwwwwsqdqsdqwdqdwqddwq")
	txpost, err := transaction.NewTxPost(u1, []byte("first"), file, singers)
	if err != nil {
		return "", err
	}
	sertxpost, err := txpost.SerializeTx()
	if err != nil {
		return "", err
	}
	return sertxpost, nil
}

func CreateTransferTx() (string, error) {
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	tx, err := transaction.NewTxTransfer(u1, []byte("first"), u2.Address(), 100)
	if err != nil {
		return "", err
	}
	sertx, err := tx.SerializeTx()
	if err != nil {
		return "", err
	}
	return sertx, nil
}

func TestDeserializeTx(t *testing.T) {
	serpost, err := CreatePostTx()
	if err != nil {
		t.Fatal(err)
	}
	tx1, err := DeserializeTx(serpost)
	if err != nil {
		t.Fatal(err)
	}
	s, err := tx1.SerializeTx()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)

}
