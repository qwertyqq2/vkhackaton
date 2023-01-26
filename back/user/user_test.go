package user

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
)

func TestParseUser(t *testing.T) {
	pk, _ := crypto.GenerateRSAPrivate()
	u := NewUser(pk)
	fmt.Println("user address", u.Addr.String())
	ustr := u.Addr.String()
	addr := ParseAddress(ustr)
	fmt.Println("copy address", addr.String())
}
