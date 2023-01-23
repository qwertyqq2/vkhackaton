package core

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/syncbc"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

type ConfBc struct {
	MinToken uint64 `json:"minTokens"`
}

func DefaultConf() *ConfBc {
	return &ConfBc{
		MinToken: 5,
	}
}

type Blockchain struct {
	conf          *ConfBc
	coll          *files.Collector
	lastBlock     uint64
	flashInterval uint64
	snap          values.Bytes
	genesisBlock  *types.Block

	sm *syncbc.SyncBcMutex
	wg sync.WaitGroup
}

func NewBlockchain(creator *user.Address) (*Blockchain, error) {
	coll, err := files.NewCollector()
	if err != nil {
		return nil, err
	}
	conf := DefaultConf()
	genesis := types.NewGenesisBLock(creator)
	snap, err := coll.State()
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		conf:          conf,
		coll:          coll,
		lastBlock:     uint64(0),
		flashInterval: uint64(0),
		snap:          snap,
		genesisBlock:  genesis,
		wg:            sync.WaitGroup{},
		sm:            syncbc.NewSyncBc(),
	}, nil
}

func (bc *Blockchain) InsertChain(blocks types.Blocks) error {
	for i := 1; i < len(blocks); i++ {
		cur, prev := blocks[i], blocks[i-1]
		if cur.Number != prev.Number+1 || !bytes.Equal(cur.PrevBlock, prev.HashBlock) ||
			!bytes.Equal(cur.PrevSnap, prev.CurShap) {
			return fmt.Errorf("incorrect base data block")
		}
	}
	ok := bc.sm.TryLock()
	if !ok {
		return fmt.Errorf("Chain is stopped")
	}
	defer bc.sm.Unlock()
	return bc.insertChain(blocks)
}

func (bc *Blockchain) insertChain(blocks types.Blocks) error {
	var (
		n       = len(blocks)
		curSnap = bc.snap
		//last    = bc.lastBlock
		dataBlocks values.Bytes
		wg         sync.WaitGroup
	)

	fs := make([]*files.File, 0)
	for i := 0; i < n; i++ {
		go func(i int) {
			wg.Add(1)
			if !bc.verifyBlock(blocks[i]) {

			}
			dataBlocks = blocks[i].Data()[i]
			fs = append(fs, files.NewFile(string(dataBlocks)))
		}(i)
	}
	wg.Wait()
	tempState := bc.coll

	return nil
}

func (bc *Blockchain) verifyBlock(block *types.Block) bool {
	return true
}
