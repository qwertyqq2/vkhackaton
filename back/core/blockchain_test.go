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
)

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

func (bc *Blockchain) printbc() error {
	blocks, err := bc.dblevel.getBlocks()
	if err != nil {
		return nil
	}
	if len(blocks) == 0 {
		return fmt.Errorf("nil size blocks")
	}
	for _, b := range blocks {
		fmt.Println(b.SerializeBlock())
	}
	return nil
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

func TestAddBlock(t *testing.T) {
	creator, pkCreator, users := initUser()
	bc := NewTestingBC(creator.Addr)
	txsTransfer := initTransferTxs(creator, pkCreator, bc.lastHashBlock, users)
	_, err := bc.AddBlock(creator, txsTransfer...)
	if err != nil {
		log.Fatal(err)
	}
	bc.printbc()
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
	b := types.NewBlock(bc.lastNumber, bc.lastHashBlock, bc.lastSnap, creator.Addr)
	b.InserTxs(txs...)
	err := b.Accept(creator)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func PostTxs(prevHash values.Bytes, users ...*user.User) []types.Transaction {
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
