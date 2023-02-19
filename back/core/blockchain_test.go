package core

import (
	"fmt"
	"log"
	"testing"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
	"github.com/stretchr/testify/assert"
)

type lastState struct {
	lastHash   values.Bytes
	lastSnap   values.Bytes
	lastNumber uint64
}

func GetLastState(h, s values.Bytes, n uint64) lastState {
	return lastState{
		lastHash:   h,
		lastSnap:   s,
		lastNumber: n,
	}
}

func AddBlock(u *user.User, ls lastState, txs ...types.Transaction) (*types.Block, error) {
	block := types.NewBlock(ls.lastNumber, ls.lastHash, ls.lastSnap, u.Addr, txs...)
	return block, nil
}

func initUser() (*user.User, *ring.PrivateKey, []*user.User) {
	users := make([]*user.User, 0)
	pk1 := ring.GeneratePrivate()
	creator := user.NewUser(pk1)
	pkCreator := pk1
	pk2 := ring.GeneratePrivate()
	users = append(users, user.NewUser(pk2))
	pk3 := ring.GeneratePrivate()
	users = append(users, user.NewUser(pk3))
	return creator, pkCreator, users
}

func postTxs(prevHash values.Bytes, users ...*user.User) []types.Transaction {
	dataFile1 := files.NewFile("is first file for me")
	dataFile2 := files.NewFile("is second file for me")
	dataFile3 := files.NewFile("is third file for me")
	dataFiles := []*files.File{dataFile1, dataFile2, dataFile3}

	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	singers := []*user.Address{u2.Addr, u3.Addr}

	txs := make([]types.Transaction, 0)
	for i, u := range users {
		tx, err := transaction.NewTxPost(u, prevHash, dataFiles[i], singers)
		if err != nil {
			log.Fatal(err)
		}
		txs = append(txs, tx)
	}
	return txs
}

func initTransferTxs(creator *user.User, pkCreator *ring.PrivateKey, prevHash values.Bytes, users []*user.User) []types.Transaction {
	txs := make([]types.Transaction, 0)
	for _, u := range users {
		tx, err := transaction.NewTxTransfer(creator, prevHash, u.Addr, 10)
		if err != nil {
			log.Fatal(err)
		}
		txs = append(txs, tx)
	}
	return txs
}

func postTransferTxs(creator *user.User, pkCreator *ring.PrivateKey, prevHash values.Bytes, users []*user.User) []types.Transaction {
	txs := make([]types.Transaction, 0)
	for _, u := range users {
		tx, err := transaction.NewTxTransfer(creator, prevHash, u.Addr, 5)
		if err != nil {
			log.Fatal(err)
		}
		txs = append(txs, tx)
	}
	return txs
}

func initTransferBlock(creator *user.User, bc *Blockchain, txs []types.Transaction) (*types.Block, error) {
	b := types.NewBlock(bc.lastNumber, bc.lastHashBlock, bc.lastSnap, creator.Addr, txs...)
	err := b.Accept(creator)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (bc *Blockchain) printbc() string {
	blocks, err := bc.dblevel.getBlocks()
	if err != nil {
		return ""
	}
	if len(blocks) == 0 {
		return ""
	}
	for _, b := range blocks {
		ser, err := b.SerializeBlock()
		if err != nil {
			return ""
		}
		fmt.Println(ser)
	}
	return ""
}

func NewTestingBC(creator *user.User) (*Blockchain, *types.Block) {
	bc, err := NewBlockchainWithGenesis(creator)
	if err != nil {
		log.Fatal(err)
	}
	return bc, bc.lastBlock
}

func TestNewBlock(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc, _ := NewTestingBC(creator)
	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
	block, err := initTransferBlock(creator, bc, txsTransfer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.SerializeBlock())
}

func TestNewBlocks(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc, _ := NewTestingBC(creator)
	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
	block, err := initTransferBlock(creator, bc, txsTransfer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.SerializeBlock())
}

func TestAddBlock(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc, _ := NewTestingBC(creator)
	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
	block, err := bc.AddBlock(creator, txsTransfer...)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("InsertBlock", func(t *testing.T) {
		err := bc.InsertChain(block)
		if err != nil {
			t.Fatal(err)
		}
	})
	bc.printbc()
}

func TestAddBlocks(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc, gen := NewTestingBC(creator)

	var (
		block1  *types.Block
		block2  *types.Block
		block3  *types.Block
		genesis *types.Block
	)

	genesis = gen

	t.Run("InsertBlocks", func(t *testing.T) {
		txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
		b1, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.InsertChain(b1); err != nil {
			t.Fatal(err)
		}

		block1 = b1

		b2, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.InsertChain(b2); err != nil {
			t.Fatal(err)
		}
		block2 = b2
		b3, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.InsertChain(b3); err != nil {
			t.Fatal(err)
		}

		block3 = b3
		assert.Equal(t, block3.CurShap, bc.lastSnap)
		assert.Equal(t, block3.HashBlock, bc.lastHashBlock)
		assert.Equal(t, block3.Number, bc.lastNumber)
		assert.Equal(t, block3.PrevBlock, block2.HashBlock)

	})

	t.Run("InsertChainWithGenesis", func(t *testing.T) {
		pk := ring.GeneratePrivate()
		client := user.NewUser(pk)
		bc1, err := NewBlockchainExternal(client.Address(), genesis, block1, block2, block3)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, block3.CurShap, bc1.lastSnap)
		assert.Equal(t, block3.HashBlock, bc1.lastHashBlock)
		assert.Equal(t, block3.Number, bc1.lastNumber)
	})

	t.Run("InsertChain", func(t *testing.T) {
		pk := ring.GeneratePrivate()
		client := user.NewUser(pk)
		bc3, err := NewBlockchainExternal(client.Address(), genesis, block1)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc3.InsertChain(block2, block3); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, block3.CurShap, bc3.lastSnap)
		assert.Equal(t, block3.HashBlock, bc3.lastHashBlock)
		assert.Equal(t, block3.Number, bc3.lastNumber)
	})

	t.Run("reorgChain", func(t *testing.T) {
		blocks1 := []*types.Block{block2, block1, block3}
		blocks2 := []*types.Block{block3, block2, block1}

		var finalBlocks = []*types.Block{block1, block2, block3}

		regBlocks := make([]*types.Block, 3)

		if needReorgBlocks(blocks1) {
			rg, err := reorgBlocks(blocks1)
			if err != nil {
				t.Fatal(err)
			}
			copy(regBlocks, rg)
		} else {
			copy(regBlocks, blocks1)
		}
		for i, b := range regBlocks {
			assert.Equal(t, b.HashBlock, finalBlocks[i].HashBlock)
			assert.Equal(t, b.Number, finalBlocks[i].Number)
			assert.Equal(t, b.CurShap, finalBlocks[i].CurShap)
		}

		if needReorgBlocks(blocks2) {
			rg, err := reorgBlocks(blocks2)
			if err != nil {
				t.Fatal(err)
			}
			copy(regBlocks, rg)
		} else {
			copy(regBlocks, blocks2)
		}
		for i, b := range regBlocks {
			assert.Equal(t, b.HashBlock, finalBlocks[i].HashBlock)
			assert.Equal(t, b.Number, finalBlocks[i].Number)
			assert.Equal(t, b.CurShap, finalBlocks[i].CurShap)
		}
	})

	t.Run("InsertReorgChain", func(t *testing.T) {
		pk := ring.GeneratePrivate()
		client := user.NewUser(pk)
		bc2, err := NewBlockchainExternal(client.Address(), genesis, block3, block2, block1)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, block3.CurShap, bc2.lastSnap)
		assert.Equal(t, block3.HashBlock, bc2.lastHashBlock)
		assert.Equal(t, block3.Number, bc2.lastNumber)
	})
}

func TestRollbackChain(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc, gen := NewTestingBC(creator)

	var (
		block1  *types.Block
		block2  *types.Block
		block3  *types.Block
		genesis *types.Block
	)

	genesis = gen

	func() {
		txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
		b1, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.InsertChain(b1); err != nil {
			t.Fatal(err)
		}

		block1 = b1

		b2, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.InsertChain(b2); err != nil {
			t.Fatal(err)
		}
		block2 = b2
		b3, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.InsertChain(b3); err != nil {
			t.Fatal(err)
		}

		block3 = b3
	}()
	t.Run("RollbackChain", func(t *testing.T) {
		pk := ring.GeneratePrivate()
		client := user.NewUser(pk)

		bc2, err := NewBlockchainExternal(client.Address(), genesis, block1)
		if err != nil {
			t.Fatal(err)
		}

		posttxsTransfer := postTransferTxs(creator, pkCreator, bc.lastHashBlock, users)

		mb, err := bc2.AddBlock(creator, posttxsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc2.InsertChain(mb); err != nil {
			t.Fatal(err)
		}

		if err := bc2.RollbackChain(block2, block3); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, block3.CurShap, bc2.lastSnap)
		assert.Equal(t, block3.HashBlock, bc2.lastHashBlock)
		assert.Equal(t, block3.Number, bc2.lastNumber)
	})
}
