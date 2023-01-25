package core

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto"
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
	dblevel       *levelDB
	lastNumber    uint64
	lastBlock     *types.Block
	lastHashBlock values.Bytes
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
	ldb, err := NewLevelDB()
	if err != nil {
		return nil, err
	}
	genesis := types.NewGenesisBLock(creator)
	snap, err := coll.Snap()
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		conf:          conf,
		coll:          coll,
		dblevel:       ldb,
		lastNumber:    uint64(0),
		flashInterval: uint64(0),
		lastHashBlock: values.Bytes("first"),
		snap:          snap,
		genesisBlock:  genesis,
		wg:            sync.WaitGroup{},
		sm:            syncbc.NewSyncBc(),
	}, nil
}

func (bc *Blockchain) InsertChain(blocks types.Blocks) error {
	if len(blocks) == 0 {
		return fmt.Errorf("nil blocks")
	}
	for i := 1; i < len(blocks); i++ {
		cur, prev := blocks[i], blocks[i-1]
		if cur.Number != prev.Number+1 || !bytes.Equal(cur.PrevBlock, prev.HashBlock) ||
			!bytes.Equal(cur.PrevSnap, prev.CurShap) {
			return fmt.Errorf("incorrect base data block")
		}
	}

	var (
		n       = len(blocks)
		errChan chan error
		wg      sync.WaitGroup
	)

	for i := 0; i < n; i++ {
		go func(i int) {
			wg.Add(1)
			if err := blocks[i].EmptyBlock(); err != nil {
				errChan <- err
			}
		}(i)
	}
	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
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
		block *types.Block
	)
	iterChain := newIterator(blocks)
	validator := newValidator(bc)
	fblock, _ := iterChain.next()

	snap, err := bc.coll.Snap()
	if err != nil {
		return err
	}

	snap, err = validator.add(snap, fblock.Transactions()...)
	if err != nil {
		return err
	}

	for ; err != nil && block != nil; block, err = iterChain.next() {
		snap, err = validator.add(snap, block.Transactions()...)
		if err != nil {
			return err
		}
	}

	iterChain.back()
	block, _ = iterChain.next()
	for ; err != nil && block != nil; block, err = iterChain.next() {
		serblock, err := block.SerializeBlock()
		if err != nil {
			return err
		}
		err = bc.dblevel.insertBlock(crypto.Base64EncodeString(block.HashBlock), serblock)
		if err != nil {
			return err
		}
		for _, tx := range block.Transactions() {
			switch tx.GetType() {
			case transaction.TypePostTx:
				err := bc.coll.InsertFile(files.NewFile(string(tx.GetData())))
				if err != nil {
					return err
				}
			case transaction.TypeTransferTx:
				err := bc.coll.SubBalance(user.ParseAddress(tx.GetSender()), tx.GetValue())
				if err != nil {
					return err
				}
				err = bc.coll.AddBalance(user.ParseAddress(tx.GetReceiver()), tx.GetValue())
				if err != nil {
					return err
				}
			}
		}
	}
	bc.lastHashBlock = block.HashBlock
	bc.lastBlock = block
	bc.lastNumber = block.Number
	bc.snap = snap
	return nil
}

func (bc *Blockchain) AddBlock(u *user.User, txs ...types.Transaction) (*types.Block, error) {
	validator := newValidator(bc)
	snap, err := bc.coll.Snap()
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {
		snap, err = validator.add(snap, tx)
		if err != nil {
			return nil, err
		}
	}
	block := types.NewBlock(bc.lastNumber, bc.lastHashBlock, bc.snap, u.Addr)
	if err := block.Pow(); err != nil {
		return nil, err
	}
	serblock, err := block.SerializeBlock()
	if err != nil {
		return nil, err
	}
	err = bc.dblevel.insertBlock(crypto.Base64EncodeString(block.HashBlock), serblock)
	if err != nil {
		return nil, err
	}
	for _, tx := range block.Transactions() {
		switch tx.GetType() {
		case transaction.TypePostTx:
			err := bc.coll.InsertFile(files.NewFile(string(tx.GetData())))
			if err != nil {
				return nil, err
			}
		case transaction.TypeTransferTx:
			err := bc.coll.SubBalance(user.ParseAddress(tx.GetSender()), tx.GetValue())
			if err != nil {
				return nil, err
			}
			err = bc.coll.AddBalance(user.ParseAddress(tx.GetReceiver()), tx.GetValue())
			if err != nil {
				return nil, err
			}
		}
	}
	bc.lastHashBlock = block.HashBlock
	bc.lastBlock = block
	bc.lastNumber = block.Number
	bc.snap = snap
	return block, nil
}
