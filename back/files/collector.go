package files

import (
	"github.com/qwertyqq2/filebc/files/state/xorstate"
)

const (
	lenHash = 32
)

type сollector struct {
	ldb *levelDB

	state State
}

func Collector() (*сollector, error) {
	l, err := LoadLevel()
	if err != nil {
		return nil, err
	}
	return &сollector{
		ldb:   l,
		state: xorstate.NewXorState(lenHash),
	}, nil
}

func (c *сollector) State(fs ...*File) ([]byte, error) {
	files, err := c.ldb.allFiles()
	if err != nil {
		return nil, err
	}
	ids := make([][]byte, 0)
	for _, f := range files {
		ids = append(ids, f.Id)
	}
	if len(fs) > 0 {
		for _, f := range fs {
			ids = append(ids, f.Id)
		}
	}
	return c.state.Get(ids...), nil
}
