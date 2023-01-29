package transaction

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
)

func TestSerializePost(t *testing.T) {
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	singers := []*user.Address{u2.Addr, u3.Addr}
	file := files.NewFile("first fileqweqweqwweqwwwwwwwwwwwwwwwwwwwwwwwwsqdqsdqwdqdwqddwq")
	tx, err := NewTxPost(u1, []byte("first"), file, singers)
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
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	pk4 := ring.GeneratePrivate()
	u4 := user.NewUser(pk4)
	pk5 := ring.GeneratePrivate()
	u5 := user.NewUser(pk5)
	pk6 := ring.GeneratePrivate()
	u6 := user.NewUser(pk6)
	singers := []*user.Address{u2.Addr, u3.Addr, u4.Addr, u5.Addr, u6.Addr}
	file := files.NewFile("first fileqweqweqwweqwwwwwwwwwwwwwwwwwwwwwwwwsqdqsdqwdqdwqddwq")
	tx, err := NewTxPost(u1, []byte("first"), file, singers)
	if err != nil {
		t.Fatal(err)
	}
	f := tx.Valid()
	if !f {
		t.Log("NOT VERIFY")
	} else {
		t.Log("VERIFY")
	}
	fmt.Println(tx.SerializeTx())
}
