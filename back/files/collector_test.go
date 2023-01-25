package files

import (
	"fmt"
	"testing"

	"github.com/qwertyqq2/filebc/crypto"
)

func TestCollector(t *testing.T) {
	c, err := NewCollector()
	if err != nil {
		t.Fatal(err)
	}
	state, err := c.Snap()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("state ", crypto.Base64EncodeString(state))

}
