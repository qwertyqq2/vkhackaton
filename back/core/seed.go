package core

import (
	"fmt"
	"sync"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/files/state/xorstate"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

type Seed interface {
	Snap() (values.Bytes, error)

	Balance(address *user.Address) (uint64, error)

	InsertFile(file *files.File) error

	AddFile(snapState values.Bytes, fs ...*files.File) values.Bytes

	AddUser(snapState values.Bytes, users ...*user.User) values.Bytes

	AddBalance(address *user.Address, delta uint64) error

	SubBalance(address *user.Address, delta uint64) error

	RemoveFile(id values.Bytes) error
}

type seed struct {
	store map[string]interface{}
	state *xorstate.XorState

	mu sync.Mutex
}

func NewSeed(snap values.Bytes) *seed {
	return &seed{
		store: map[string]interface{}{},
		state: xorstate.NewXorState(),
	}
}

func (s *seed) Snap() (values.Bytes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := make([]values.Bytes, 0)
	for k := range s.store {
		ids = append(ids, crypto.Base64DecodeString(k))
	}
	return s.state.Get(ids...), nil
}

func (s *seed) InsertFile(f *files.File) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.store[crypto.Base64EncodeString(f.Id)]; ok {
		return fmt.Errorf("file already exist")
	}
	s.store[crypto.Base64EncodeString(f.Id)] = f
	return nil
}

func (s *seed) AddBalance(address *user.Address, delta uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.store[address.String()].(*user.User)
	if !ok {
		return fmt.Errorf("not exist such user")
	}
	u.Balance += delta
	s.store[address.String()] = u
	return nil
}

func (s *seed) SubBalance(address *user.Address, delta uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.store[address.String()].(*user.User)
	if !ok {
		return fmt.Errorf("not exist such user")
	}
	u.Balance -= delta
	s.store[address.String()] = u
	return nil
}

func (s *seed) Balance(address *user.Address) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.store[address.String()].(*user.User)
	if !ok {
		return 0, fmt.Errorf("not exist such user")
	}
	return u.Balance, nil
}

func (s *seed) RemoveFile(id values.Bytes) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, crypto.Base64EncodeString(id))
	return nil
}

func (s *seed) AddUser(snapState values.Bytes, users ...*user.User) values.Bytes {
	for _, u := range users {
		snapState = s.state.Add(snapState, u.Hash())
	}
	return snapState
}

func (s *seed) AddFile(snapState values.Bytes, fs ...*files.File) values.Bytes {
	for _, f := range fs {
		snapState = s.state.Add(snapState, f.Id)
	}
	return snapState
}
