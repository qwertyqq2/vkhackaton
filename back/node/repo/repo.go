package reponode

import (
	"github.com/qwertyqq2/filebc/crypto/ring"
)

type Repo interface {
	PrivateKey() (*ring.PrivateKey, error)

	PeerStorePath() (string, error)

	SaveKey(key string) error
}
