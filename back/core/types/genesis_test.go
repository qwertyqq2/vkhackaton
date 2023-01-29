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
	gen := NewGenesisBLock(u1.Address())
	fmt.Println(gen)
}
