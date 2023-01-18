package types

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/user"
)

func TestGenesis(t *testing.T) {
	pk1, err := crypto.GenerateRSAPrivate()
	if err != nil {
		t.Fatal(err)
	}
	u1 := user.NewUser(pk1)
	gen := NewGenesisBLock(u1.Address())
	fmt.Println(gen)
}
