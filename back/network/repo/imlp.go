package repo

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/libp2p/go-libp2p/core/crypto"
)

var (
	ErrRepoNotExist = errors.New("repo not exist err")
)

type repositoryImpl struct {
	sync.RWMutex
	rootpath string
}

func Open(path string) (Repo, error) {
	if !exist(path) {
		if err := os.Mkdir(path, 0777); err != nil {
			return nil, err
		}
	}
	return &repositoryImpl{
		rootpath: path,
	}, nil
}

func exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func generatePriv() (crypto.PrivKey, error) {
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func (r *repositoryImpl) PrivateKey() (crypto.PrivKey, error) {
	r.Lock()
	defer r.Unlock()
	if !exist(r.rootpath) {
		return nil, ErrRepoNotExist
	}
	path := filepath.Join(r.rootpath, "peer.key")
	if !exist(path) {
		pk, err := generatePriv()
		if err != nil {
			return nil, err
		}
		bytesPriv, err := crypto.MarshalPrivateKey(pk)
		if err != nil {
			return nil, err
		}
		if err := r.SaveKey(string(bytesPriv)); err != nil {
			return nil, err
		}
		return pk, nil
	}
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(buff)
}

func (r *repositoryImpl) PeerStorePath() (string, error) {
	r.Lock()
	defer r.Unlock()
	path := filepath.Join(r.rootpath, "peer", "storage")
	if !exist(r.rootpath) {
		return "", ErrRepoNotExist
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		return "", err
	}
	return path, nil
}

func (r *repositoryImpl) SaveKey(key string) error {
	path := filepath.Join(r.rootpath, "peer.key")
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(key), 0777)
}
