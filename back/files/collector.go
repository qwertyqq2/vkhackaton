package files

import (
	"github.com/qwertyqq2/filebc/files/state/xorstate"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

const (
	lenHash = 32
)

type InState struct {
	addr *user.Address
	bal  uint64
	file *File
}

type Collector struct {
	ldb *levelDB

	state State
}

func NewCollector() (*Collector, error) {
	l, err := NewLevelDB()
	if err != nil {
		return nil, err
	}
	return &Collector{
		ldb:   l,
		state: xorstate.NewXorState(lenHash),
	}, nil
}

func LoadCollector() (*Collector, error) {
	l, err := LoadLevel()
	if err != nil {
		return nil, err
	}
	return &Collector{
		ldb:   l,
		state: xorstate.NewXorState(lenHash),
	}, nil
}

func (c *Collector) State(fs ...*File) (values.Bytes, error) {
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
	return c.state.Get(ids...), nil
}

func (c *Collector) Balance(address *user.Address) (uint64, error) {
	return c.ldb.getBalance(address.String())
}

func (c *Collector) AddFile(file *File) error {
	return c.ldb.insertFile(file)
}

func (c *Collector) AddBalance(address *user.Address, delta uint64) error {
	return c.ldb.addBalance(address.String(), delta)
}
