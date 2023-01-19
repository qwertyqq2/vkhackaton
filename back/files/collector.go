package files

import (
	"github.com/qwertyqq2/filebc/files/state/xorstate"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

const (
	lenHash = 32
)

type Collector struct {
	ldb *levelDB

	state State
}

func NewCollector() (*Collector, error) {
	l, err := LoadLevel()
	if err != nil {
		return nil, err
	}
	return &Collector{
		ldb:   l,
		state: xorstate.NewXorState(lenHash),
	}, nil
}

func (c *Collector) State(fs ...*File) ([]byte, error) {
	files, err := c.ldb.allFiles()
	if err != nil {
		return nil, err
	}
	ids := make([]values.Bytes, 0)
	for _, f := range files {
		ids = append(ids, f.Id)
	}
	if len(fs) > 0 {
		for _, f := range fs {
			ids = append(ids, f.Id)
		}
	}
	usersWrap, err := c.ldb.getUsers()
	if err != nil {
		return nil, err
	}
	if len(usersWrap) > 0 {
		users := make([]*user.User, 0)
		for _, uw := range usersWrap {
			addr := user.ParseAddress(uw.Addr)
			users = append(users, &user.User{
				Addr:    addr,
				Balance: uint64(uw.Bal),
			})
		}
		for _, u := range users {
			ids = append(ids, u.Hash())
		}
	}
	return c.state.Get(ids...), nil
}
