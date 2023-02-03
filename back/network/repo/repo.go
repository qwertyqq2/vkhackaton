package repo

import (
	"github.com/libp2p/go-libp2p/core/crypto"
)

type Repo interface {
	PrivateKey() (crypto.PrivKey, error)

	PeerStorePath() (string, error)

	SaveKey(key string) error
}
