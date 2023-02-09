package files

import (
	"fmt"
	"sync"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files/state/xorstate"
	"github.com/qwertyqq2/filebc/syncbc"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

type Collector struct {
	syncbc.SyncBcMutex
	wg sync.WaitGroup

	snap uint64

	ldb   *levelDB
	state State
}

func NewCollector() (*Collector, error) {
	l, err := NewLevelDB()
	if err != nil {
		return nil, err
	}
	return &Collector{
		ldb:         l,
		state:       xorstate.NewXorState(),
		SyncBcMutex: *syncbc.NewSyncBc(),
	}, nil
}

func (c *Collector) Snap() (values.Bytes, error) {
	var (
		exitChan = make(chan bool)
		ids      = make([]values.Bytes, 0)
	)
	c.wg.Add(2)
	go func() {
		defer c.wg.Done()
		c.Locking()
		files, err := c.ldb.allFiles()
		c.Unlock()
		if err != nil {
			exitChan <- true
			return
		}
		for _, f := range files {
			ids = append(ids, f.Id)
		}
	}()
	go func() {
		defer c.wg.Done()
		c.Locking()
		usersWrap, err := c.ldb.getUsers()
		c.Unlock()
		if err != nil {
			exitChan <- true
			return
		}
		for i := 0; i < len(usersWrap); i++ {
			addr, err := user.ParseAddress(usersWrap[i].Addr)
			if err != nil {
				exitChan <- true
				return
			}
			u := &user.User{
				Addr:    addr,
				Balance: uint64(usersWrap[i].Bal),
			}
			ids = append(ids, u.Hash())
		}
	}()
	c.wg.Wait()
	select {
	case <-exitChan:
		return nil, fmt.Errorf("err parsing")
	default:
	}
	return c.state.Get(ids...), nil
}

func (c *Collector) AddUser(snapState values.Bytes, users ...*user.User) values.Bytes {
	for _, u := range users {
		snapState = c.state.Add(snapState, u.Hash())
	}
	return snapState
}

func (c *Collector) AddFile(snapState values.Bytes, fs ...*File) values.Bytes {
	for _, f := range fs {
		snapState = c.state.Add(snapState, f.Id)
	}
	return snapState
}

func (c *Collector) Balance(address *user.Address) (uint64, error) {
	return c.ldb.getBalance(address.String())
}

func (c *Collector) InsertFile(file *File) error {
	return c.ldb.insertFile(file)
}

func (c *Collector) AddBalance(address *user.Address, delta uint64) error {
	return c.ldb.addBalance(address.String(), delta)
}

func (c *Collector) SubBalance(address *user.Address, delta uint64) error {
	return c.ldb.subBalance(address.String(), delta)
}

func (c *Collector) RemoveFile(id values.Bytes) error {
	return c.ldb.removeFileById(crypto.Base64EncodeString(id))
}
