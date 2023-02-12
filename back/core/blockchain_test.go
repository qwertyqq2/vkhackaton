package core

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
	"github.com/stretchr/testify/assert"
)

func (bc *Blockchain) AddBlock(u *user.User, txs ...types.Transaction) (*types.Block, error) {
	validator := newValidator(bc.coll)
	snap, err := bc.coll.Snap()
	if err != nil {
		return nil, err
	}

	// for _, tx := range txs {
	// 	snap, err = validator.add(snap, tx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	snap, err = validator.add(snap, txs...)
	if err != nil {
		return nil, err
	}
	block := types.NewBlock(bc.lastNumber, bc.lastHashBlock, bc.lastSnap, u.Addr, txs...)
	if err := block.Accept(u); err != nil {
		return nil, err
	}
	block.CurShap = snap

	block.Time = time.Now().Format(time.RFC3339)
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

func NewTestingBC(creator *user.Address) *Blockchain {
	bc, err := NewBlockchain(creator)
	if err != nil {
		log.Fatal(err)
	}
	return bc
}

func TestNewBlock(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc := NewTestingBC(creator.Addr)
	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
	block, err := initTransferBlock(creator, bc, txsTransfer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.SerializeBlock())
}

func TestNewBlocks(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc := NewTestingBC(creator.Addr)
	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
	block, err := initTransferBlock(creator, bc, txsTransfer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.SerializeBlock())
}

// func TestAddBlock(t *testing.T) {
// 	creator, pkCreator, users := initUser()
// 	bc := NewTestingBC(creator.Addr)
// 	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
// 	block, err := bc.AddBlock(creator, txsTransfer...)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Run("InsertBlock", func(t *testing.T) {
// 		err := bc.InsertChain(block)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})
// 	bc.printbc()
// }

func TestAddBlocks(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc := NewTestingBC(creator.Addr)

	t.Run("InsertBlocks", func(t *testing.T) {
		txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
		block1, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		err = bc.InsertChain(block1)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(crypto.Base64EncodeString(block1.CurShap), crypto.Base64EncodeString(bc.lastSnap))
		fmt.Println(bytes.Equal(block1.CurShap, bc.lastSnap))
		block2, err := bc.AddBlock(creator, txsTransfer...)
		if err != nil {
			t.Fatal(err)
		}
		err = bc.InsertChain(block2)
		if err != nil {
			t.Fatal(err)
		}
		// val := assert.Equal(t, block2.CurShap, bc.lastSnap)
		// fmt.Println(val)
		fmt.Println(crypto.Base64EncodeString(block2.CurShap), crypto.Base64EncodeString(bc.lastSnap))
		fmt.Println(bytes.Equal(block2.CurShap, bc.lastSnap))
		assert.Equal(t, block2.HashBlock, bc.lastHashBlock)
		assert.Equal(t, block2.Number, bc.lastNumber)

	})
}
