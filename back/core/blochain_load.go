package core

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/syncbc"
	"github.com/qwertyqq2/filebc/user"
)

var (
	ErrExistBc = errors.New("exist bc err")
	ErrLoadBc  = errors.New("load bc err")
)

func ExistBc(uname string) (bool, error) {
	_, err1 := files.LoadCollector(uname)
	_, err2 := loadLevel(uname)
	if err1 != nil && err2 != nil {
		return false, nil
	}
	if err1 == nil && err2 == nil {
		return true, nil
	}
	if err1 != nil && err2 == nil || err1 == nil && err2 != nil {
		return false, ErrExistBc
	}
	return false, nil
}

func LoadBlockchain(addr *user.Address) (*Blockchain, error) {
	var uname = addr.String()[:2]
	ex, err := ExistBc(uname)
	if err != nil {
		return nil, fmt.Errorf("you have corrupted data")
	}
	if !ex {
		return nil, ErrLoadBc
	}
	coll, err := files.LoadCollector(uname)
	if err != nil {
		return nil, err
	}
	l, err := loadLevel(uname)
	if err != nil {
		return nil, err
	}
	snap, err := coll.Snap()
	if err != nil {
		return nil, err
	}
	lastBlock, err := l.lastBlock()
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		conf:          DefaultConf(),
		coll:          coll,
		dblevel:       l,
		lastNumber:    lastBlock.Number,
		flashInterval: uint64(0),
		lastHashBlock: lastBlock.HashBlock,
		lastSnap:      snap,
		wg:            sync.WaitGroup{},
		sm:            syncbc.NewSyncBc(),
	}, nil
}

func NewBlockchainExternal(client *user.Address, blocks ...*types.Block) (*Blockchain, error) {
	var regBlocks = make(types.Blocks, len(blocks))

	if needReorgBlocks(blocks) {
		rg, err := reorgBlocks(blocks)
		if err != nil {
			return nil, err
		}
		copy(regBlocks, rg)
	} else {
		copy(regBlocks, blocks)
	}

	gen := regBlocks[0]
	addrCreator, err := user.ParseAddress(gen.Miner)
	if err != nil {
		return nil, err
	}

	coll, err := files.NewCollector(client.String()[:2])
	if err != nil {
		return nil, err
	}
	conf := DefaultConf()
	ldb, err := NewLevelDB(client.String()[:2])
	if err != nil {
		return nil, err
	}
	valGen := gen.Transactions()[0].GetValue()
	if err := coll.AddBalance(addrCreator, valGen); err != nil {
		return nil, err
	}
	curSnap, err := coll.Snap()
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(gen.CurShap, curSnap) || !bytes.Equal(gen.PrevBlock, []byte(types.GenPrevblock)) || !bytes.Equal(gen.PrevSnap, []byte(types.GenPrevSnap)) {
		return nil, fmt.Errorf("gen valid err")
	}

	if !gen.Valid() {
		return nil, fmt.Errorf("gen hash valid err")
	}
	bc := &Blockchain{
		conf:          conf,
		coll:          coll,
		dblevel:       ldb,
		lastNumber:    uint64(1),
		flashInterval: uint64(0),
		lastHashBlock: gen.HashBlock,
		lastSnap:      curSnap,
		genesisBlock:  gen,
		wg:            sync.WaitGroup{},
		sm:            syncbc.NewSyncBc(),
	}

	sergenblock, err := gen.SerializeBlock()
	if err != nil {
		return nil, err
	}
	if err := bc.dblevel.insertBlock(crypto.Base64EncodeString(gen.HashBlock), sergenblock); err != nil {
		return nil, err
	}
	bc.lastNumber = gen.Number
	bc.lastBlock = gen
	bc.lastHashBlock = gen.HashBlock
	if len(regBlocks[1:]) == 0 {
		return bc, nil
	}
	if err := bc.InsertChain(regBlocks[1:]...); err != nil {
		return nil, err
	}
	return bc, nil
}

func (bc *Blockchain) ReadChain() ([]string, error) {
	bs, err := bc.dblevel.getBlocks()
	if err != nil {
		return nil, err
	}
	bsstr := make([]string, len(bs))
	for _, b := range bs {
		bss, err := b.SerializeBlock()
		if err != nil {
			return nil, err
		}
		bsstr = append(bsstr, bss)
	}
	return bsstr, nil
}

func (bc *Blockchain) ReadCollUsers() ([]string, error) {
	users, err := bc.coll.LDB().GetUsers()
	if err != nil {
		return nil, err
	}
	uss := make([]string, len(users))
	for _, uw := range users {
		addr, err := user.ParseAddress(uw.Addr)
		if err != nil {
			return nil, err
		}
		u := user.GetUser(addr, uint64(uw.Bal))
		us, err := u.Serialize()
		if err != nil {
			return nil, err
		}
		uss = append(uss, us)
	}
	return uss, nil
}

func (bc *Blockchain) ReadCollFiles() ([]string, error) {
	files, err := bc.coll.LDB().GetFiles()
	if err != nil {
		return nil, err
	}
	return files, nil
}
