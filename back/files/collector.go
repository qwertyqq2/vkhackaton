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

func (c *Collector) Snap() (values.Bytes, error) {
	files, err := c.ldb.allFiles()
	if err != nil {
		return nil, err
	}
	ids := make([]values.Bytes, 0)
	for _, f := range files {
		ids = append(ids, f.Id)
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
