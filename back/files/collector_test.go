package files

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
)

func TestCollector(t *testing.T) {
	c, err := Collector()
	if err != nil {
		t.Fatal(err)
	}
	state, err := c.State()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("state ", crypto.Base64EncodeString(state))
	f1 := NewFile("example")
	state, err = c.State(f1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("state ", crypto.Base64EncodeString(state))

}
