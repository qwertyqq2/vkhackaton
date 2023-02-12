package files

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/user"
	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	c, err := NewCollector()
	if err != nil {
		t.Fatal(err)
	}
	pk1 := ring.GeneratePrivate()
	u1 := user.NewUser(pk1)
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	if err := c.AddBalance(u1.Addr, 100); err != nil {
		t.Fatal(err)
	}
	if err := c.AddBalance(u2.Addr, 200); err != nil {
		t.Fatal(err)
	}
	if err := c.AddBalance(u3.Addr, 300); err != nil {
		t.Fatal(err)
	}
	state, err := c.Snap()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("state init", crypto.Base64EncodeString(state))
	t.Run("testing iter coll", func(t *testing.T) {

		inv1 := c.State().Inverse(user.GetUser(u1.Addr, 100).Hash())
		state = c.State().Add(state, inv1)
		state = c.AddUser(state, user.GetUser(u1.Addr, 50))

		inv2 := c.State().Inverse(user.GetUser(u2.Addr, 200).Hash())
		state = c.State().Add(state, inv2)
		state = c.AddUser(state, user.GetUser(u2.Addr, 250))

		if err := c.SubBalance(u1.Addr, 50); err != nil {
			t.Fatal(err)
		}
		if err := c.AddBalance(u2.Addr, 50); err != nil {
			t.Fatal(err)
		}
		statenew, err := c.Snap()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, state, statenew)

	})
}
