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
	MinToken   uint64 `json:"minTokens"`
	InitTokens uint64 `json:"initTokens"`
	MaxBlocks  uint
}

func DefaultConf() *ConfBc {
	return &ConfBc{
		MinToken:   5,
		InitTokens: 100,
	}
}

type Blockchain struct {
	conf    *ConfBc
	coll    *files.Collector
	dblevel *levelDB

	lastNumber    uint64
	lastBlock     *types.Block
	lastHashBlock values.Bytes
	flashInterval uint64
	lastSnap      values.Bytes

	genesisBlock *types.Block

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
	err = coll.AddBalance(creator, conf.InitTokens)
	if err != nil {
		return nil, err
	}
	snap, err := coll.Snap()
	if err != nil {
		return nil, err
	}
	genesis.CurShap = snap
	bc := &Blockchain{
		conf:          conf,
		coll:          coll,
		dblevel:       ldb,
		lastNumber:    uint64(0),
		flashInterval: uint64(0),
		lastHashBlock: values.Bytes("first"),
		lastSnap:      snap,
		genesisBlock:  genesis,
		wg:            sync.WaitGroup{},
		sm:            syncbc.NewSyncBc(),
	}
	sergenblock, err := genesis.SerializeBlock()
	if err != nil {
		return nil, err
	}
	if err := bc.dblevel.insertBlock(crypto.Base64EncodeString(genesis.HashBlock), sergenblock); err != nil {
		return nil, err
	}
	bc.lastNumber = 1
	bc.lastBlock = genesis
	bc.lastHashBlock = genesis.HashBlock
	return bc, nil
}

func (bc *Blockchain) loadLastState() error {
	lastblock, err := bc.dblevel.lastBlock()
	if err != nil {
		return err
	}
	state, err := bc.coll.Snap()
	if err != nil {
		return err
	}
	bc.lastBlock = lastblock
	bc.lastHashBlock = lastblock.HashBlock
	bc.lastNumber = lastblock.Number
	bc.lastSnap = state
	return nil
}

func needReorgBlocks(blocks types.Blocks) bool {
	for i := 1; i < len(blocks); i++ {
		if blocks[i].Number-blocks[i-1].Number < 0 {
			return true
		}
		if !bytes.Equal(blocks[i].PrevBlock, blocks[i-1].HashBlock) {
			return true
		}
		if !bytes.Equal(blocks[i].PrevSnap, blocks[i-1].CurShap) {
			return true
		}
	}
	return false
}

func reorgBlocks(blocks types.Blocks) (types.Blocks, error) {
	var (
		resBlocks = make([]*types.Block, len(blocks))
		it        = 0
		n         = len(blocks)
		idx       int
	)
	for i, b := range blocks {
		if b.Number >= blocks[idx].Number {
			idx = i
		}
	}
	resBlocks[n-it-1] = blocks[idx]
	it++
	for it < len(blocks) {
		num := 0
		for i, b := range blocks {
			if b.Number < blocks[num].Number && blocks[num].Number < blocks[idx].Number {
				num = i
			}
		}
		if blocks[idx].Number-blocks[num].Number > 1 {
			return nil, fmt.Errorf("incorrect numbers in given chain")
		}
		resBlocks[n-it-1] = blocks[num]
		it++
		idx = num
	}
	for i := 1; i < len(resBlocks); i++ {
		if !bytes.Equal(resBlocks[i].PrevBlock, resBlocks[i-1].HashBlock) {
			return nil, fmt.Errorf("incorrect hashes in given chain")
		}
		if !bytes.Equal(resBlocks[i].PrevSnap, resBlocks[i-1].CurShap) {
			return nil, fmt.Errorf("incorrect snaps in given chain")
		}
	}
	return resBlocks, nil
}

func (bc *Blockchain) insertTxsLevelDb(errChan chan error, txs ...types.Transaction) error {
	for _, tx := range txs {
		bc.wg.Add(1)
		go func(tx types.Transaction) {
			defer bc.wg.Done()
			switch tx.GetType() {
			case transaction.TypePostTx:
				err := bc.coll.InsertFile(files.NewFile(string(tx.GetData())))
				if err != nil {
					errChan <- err
					return
				}
			case transaction.TypeTransferTx:
				sender, err := user.ParseAddress(tx.GetSender())
				if err != nil {
					errChan <- err
					return
				}
				err = bc.coll.SubBalance(sender, tx.GetValue())
				if err != nil {
					errChan <- err
					return
				}
				receiver, err := user.ParseAddress(tx.GetReceiver())
				if err != nil {
					errChan <- err
					return
				}
				err = bc.coll.AddBalance(receiver, tx.GetValue())
				if err != nil {
					errChan <- err
					return
				}
			}
		}(tx)
	}
	bc.wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (bc *Blockchain) InsertChain(blocks ...*types.Block) error {
	if len(blocks) == 0 {
		return fmt.Errorf("nil blocks")
	}

	var regBlocks = make(types.Blocks, len(blocks))

	if needReorgBlocks(blocks) {
		rg, err := reorgBlocks(blocks)
		if err != nil {
			return err
		}
		copy(regBlocks, rg)
	} else {
		copy(regBlocks, blocks)
	}

	if !bytes.Equal(regBlocks[0].PrevBlock, bc.lastHashBlock) || regBlocks[0].Number != bc.lastNumber+1 {
		return fmt.Errorf("its not next block")
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
		wg.Add(1)
		go func(i int) {
			if err := blocks[i].EmptyBlock(); err != nil {
				wg.Done()
				errChan <- err
			}
			wg.Done()
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
	return bc.insertChain(blocks...)
}

func (bc *Blockchain) insertChain(blocks ...*types.Block) error {
	var (
		block   *types.Block
		errChan = make(chan error)
	)

	iterChain := newIterator(blocks)
	validator := newValidator(bc.coll)
	fblock, _ := iterChain.next()

	snap, err := bc.coll.Snap()
	if err != nil {
		return err
	}
	fmt.Println("last snap", crypto.Base64EncodeString(snap))
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
	fmt.Println("new snap", crypto.Base64EncodeString(snap))

	iterChain.back()
	block, _ = iterChain.next()

	serblock, err := block.SerializeBlock()
	if err != nil {
		return err
	}
	err = bc.dblevel.insertBlock(crypto.Base64EncodeString(block.HashBlock), serblock)
	if err != nil {
		return err
	}
	if err := bc.insertTxsLevelDb(errChan, block.Transactions()...); err != nil {
		return err
	}

	for ; err != nil && block != nil; block, err = iterChain.next() {
		serblock, err := block.SerializeBlock()
		if err != nil {
			return err
		}
		err = bc.dblevel.insertBlock(crypto.Base64EncodeString(block.HashBlock), serblock)
		if err != nil {
			return err
		}
		if err := bc.insertTxsLevelDb(errChan, block.Transactions()...); err != nil {
			return err
		}
	}
	return bc.loadLastState()
}

func (bc *Blockchain) insertRewindChain(seed Seed, txs ...types.Transaction) error {
	for _, tx := range txs {
		switch tx.GetType() {
		case transaction.TypeTransferTx:
			sender, err := user.ParseAddress(tx.GetSender())
			if err != nil {
				return err
			}
			err = seed.AddBalance(sender, tx.GetValue())
			if err != nil {
				return err
			}
			receiver, err := user.ParseAddress(tx.GetReceiver())
			if err != nil {
				return err
			}
			err = seed.SubBalance(receiver, tx.GetValue())
			if err != nil {
				return err
			}

		case transaction.TypePostTx:
			file, err := files.Deserialize(string(tx.GetData()))
			if err != nil {
				return err
			}
			if err := seed.RemoveFile(file.Id); err != nil {
				return err
			}

		}
	}
	return nil
}

func (bc *Blockchain) rewindChain(idx uint64) error {
	bc.sm.Lock()
	defer bc.sm.Unlock()

	var (
		lastId = bc.lastNumber
		i      uint64
	)

	if lastId < idx {
		return fmt.Errorf("id not less last")
	}

	for i = lastId; i >= idx; i++ {
		block, err := bc.dblevel.blockById(i)
		if err != nil {
			return err
		}
		if err := bc.insertRewindChain(bc.coll, block.Transactions()...); err != nil {
			return err
		}
	}
	return bc.loadLastState()
}

func (bc *Blockchain) needRollbackChain(blocks types.Blocks) bool {
	if bc.lastNumber >= blocks[0].Number {
		return false
	}
	if bc.lastNumber-blocks[0].Number > uint64(len(blocks)) {
		return false
	}
	b, err := bc.dblevel.blockById(blocks[0].Number)
	if err != nil {
		return false
	}
	if !bytes.Equal(b.HashBlock, blocks[0].HashBlock) {
		return false
	}
	if !bytes.Equal(b.PrevBlock, blocks[0].PrevBlock) {
		return false
	}
	if !bytes.Equal(b.CurShap, blocks[0].CurShap) {
		return false
	}
	return true
}

func (bc *Blockchain) syntheticRewindChain(idx uint64, seed Seed) error {
	bc.sm.Lock()
	defer bc.sm.Unlock()
	var (
		lastId = bc.lastNumber
		i      uint64
	)

	for i = lastId; i >= idx; i++ {
		block, err := bc.dblevel.blockById(i)
		if err != nil {
			return err
		}
		if err := bc.insertRewindChain(seed, block.Transactions()...); err != nil {
			return err
		}

	}
	return nil
}

func (bc *Blockchain) RollbackChain(blocks types.Blocks) error {
	var (
		rgBlocks = make([]*types.Block, len(blocks))
		seed     = NewSeed(bc)
	)
	if needReorgBlocks(blocks) {
		rb, err := reorgBlocks(blocks)
		if err != nil {
			return err
		}
		_ = copy(rgBlocks, rb)
	} else {
		copy(rgBlocks, blocks)
	}
	if !bc.needRollbackChain(rgBlocks) {
		return fmt.Errorf("rollback dont need")
	}
	if err := bc.syntheticRewindChain(rgBlocks[0].Number, seed); err != nil {
		return fmt.Errorf("synth rewind err")
	}

	snap := rgBlocks[0].CurShap

	iterChain := newIterator(blocks)
	validator := newValidator(bc.coll)
	block, _ := iterChain.next()

	snap, err := bc.coll.Snap()
	if err != nil {
		return err
	}

	snap, err = validator.add(snap, block.Transactions()...)
	if err != nil {
		return err
	}
	for ; err != nil && block != nil; block, err = iterChain.next() {
		snap, err = validator.add(snap, block.Transactions()...)
		if err != nil {
			return err
		}
	}

	if err := bc.rewindChain(blocks[0].Number); err != nil {
		return fmt.Errorf("rewind chain err")
	}

	ok := bc.sm.TryLock()
	if !ok {
		return fmt.Errorf("Chain is stopped")
	}
	defer bc.sm.Unlock()
	return bc.insertChain(blocks...)

}
