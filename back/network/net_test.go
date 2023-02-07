package network

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/multiformats/go-multiaddr"
	"github.com/qwertyqq2/filebc/network/repo"
)

func TestMultiaddr(t *testing.T) {
	m1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/1234")
	fmt.Println(m1.String())
}

func createpk() error {
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return err
	}
	repo, err := repo.Open("node")
	if err != nil {
		return err
	}
	bytesPriv, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return err
	}
	if err := repo.SaveKey(string(bytesPriv)); err != nil {
		return err
	}
	return nil
}

func TestInitNode(t *testing.T) {
	for {
		select {
		default:
			break
		}
		break
	}
	fmt.Println("qw")
}
