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

	State() files.State

	Add(values.Bytes, ...values.Bytes) values.Bytes

	Balance(address *user.Address) (uint64, bool, error)

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

	bc   *Blockchain
	snap values.Bytes

	mu sync.Mutex
}

func NewSeed(bc *Blockchain) *seed {
	return &seed{
		store: map[string]interface{}{},
		state: xorstate.NewXorState(),
		bc:    bc,
	}
}

func (s *seed) Snap() (values.Bytes, error) {
	if s.snap == nil {
		return nil, fmt.Errorf("nil snap seed")
	}
	return s.snap, nil
}

func (s *seed) State() files.State {
	return nil
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
		bal, _, err := s.balance(address)
		if err != nil {
			return err
		}
		u := &user.User{
			Addr:    address,
			Balance: bal + delta,
		}
		s.store[address.String()] = u
		return nil
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
		bal, _, err := s.balance(address)
		if err != nil {
			return err
		}
		u := &user.User{
			Addr:    address,
			Balance: bal - delta,
		}
		s.store[address.String()] = u
	}
	u.Balance -= delta
	s.store[address.String()] = u
	return nil
}

func (s *seed) balance(address *user.Address) (uint64, bool, error) {
	bal, _, err := s.bc.coll.Balance(address)
	if err != nil {
		return 0, false, nil
	}
	s.store[address.String()] = &user.User{
		Addr:    address,
		Balance: bal,
	}
	return bal, true, nil
}

func (s *seed) Balance(address *user.Address) (uint64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.store[address.String()].(*user.User)
	if !ok {
		return s.balance(address)
	}
	return u.Balance, true, nil
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

func (s *seed) Add(snapState values.Bytes, dn ...values.Bytes) values.Bytes {
	for _, d := range dn {
		snapState = s.state.Add(snapState, d)
	}
	return snapState
}
