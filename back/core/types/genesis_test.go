package types

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/user"
)

func TestGenesis(t *testing.T) {
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	gen := NewGenesisBLock(u1.Address(), 100)
	if err := gen.AcceptGenesis(u1); err != nil {
		t.Fatal(err)
	}
	ser, err := gen.SerializeBlock()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ser)
}
