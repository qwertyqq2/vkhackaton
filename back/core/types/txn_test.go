package types

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
)

func CreatePostTx() (string, error) {
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		return "", err
	}
	u1 := user.NewUser(pk1)
	file := files.GenerateFile("first fileqweqweqwweqwwwwwwwwwwwwwwwwwwwwwwwwsqdqsdqwdqdwqddwq")
	txpost, err := transaction.NewTxPost(u1, []byte("first"), file)
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
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		return "", err
	}
	u1 := user.NewUser(pk1)
	pk2, err := crypto.GenerateRSAPrivate()
	if err != nil {
		return "", err
	}
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
	fmt.Println(tx1)

}
