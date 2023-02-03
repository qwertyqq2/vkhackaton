package user

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto/ring"
)

func TestParseUser(t *testing.T) {
	pk := ring.GeneratePrivate()
	u := NewUser(pk)
	ustr := u.Addr.String()
	addr, err := ParseAddress(ustr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("copy address", addr.String())
}

func TestShuffle(t *testing.T) {
	a := []string{"aaa", "bbb", "ccc", "ddd", "fff"}
	Shuffle(a)
	fmt.Println(a)
}
