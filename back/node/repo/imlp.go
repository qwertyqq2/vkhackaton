package reponode

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/qwertyqq2/filebc/crypto/ring"
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

func (r *repositoryImpl) PrivateKey() (*ring.PrivateKey, error) {
	r.Lock()
	defer r.Unlock()
	if !exist(r.rootpath) {
		return nil, ErrRepoNotExist
	}
	path := filepath.Join(r.rootpath, "peer.key")
	if !exist(path) {
		pk := ring.GeneratePrivate()
		bytesPriv := pk.Marshal()
		if err := r.SaveKey(string(bytesPriv)); err != nil {
			return nil, err
		}
		return pk, nil
	}
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ring.UnmarshalPrivate(buff), nil
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
