package main

import (
	"fmt"
	"log"
	"os"

	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

type node struct {
	user *user.User
	bc   *core.Blockchain
}

func NewNode() (*node, error) {
	u, _ := initUser()
	bc, err := newBc(u)
	if err != nil {
		return nil, err
	}
	return &node{
		bc:   bc,
		user: u,
	}, nil
}

func (n *node) Get() ([]string, error) {
	fscoll, err := n.bc.ReadCollFiles()
	if err != nil {
		return nil, err
	}
	fs := make([]string, len(fscoll))
	for _, f := range fscoll {
		fstr, err := files.Deserialize(f)
		if err != nil {
			return nil, err
		}
		fs = append(fs, string(fstr.Data))
	}
	return fs, nil
}

func (n *node) Set(fs ...string) error {
	dataFiles := make([]*files.File, 0)
	for _, fstr := range fs {
		f := files.NewFile(fstr)
		dataFiles = append(dataFiles, f)
	}
	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	singers := []*user.Address{u2.Addr, u3.Addr}

	txs := make([]types.Transaction, 0)
	for _, d := range dataFiles {
		tx, err := transaction.NewTxPost(n.user, n.bc.State().LastHashBlock, d, singers)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}
	block, err := n.bc.AddBlock(n.user, txs...)
	if err != nil {
		return err
	}
	if err := n.bc.InsertChain(block); err != nil {
		return err
	}
	return nil
}

func initUser() (*user.User, *ring.PrivateKey) {
	pk1 := ring.GeneratePrivate()
	creator := user.NewUser(pk1)
	return creator, pk1
}

func newBc(creator *user.User) (*core.Blockchain, error) {
	bc, err := core.NewBlockchainWithGenesis(creator)
	if err != nil {
		return nil, err
	}
	txs, err := postTxs(bc.State().LastHashBlock, creator)
	if err != nil {
		return nil, err
	}
	block, err := bc.AddBlock(creator, txs...)
	if err != nil {
		return nil, err
	}
	log.Println("add block comp")
	if err := bc.InsertChain(block); err != nil {
		return nil, err
	}
	return bc, nil
}

func postTxs(prevHash values.Bytes, u *user.User) ([]types.Transaction, error) {
	data1, err := os.ReadFile("htmlfiles/htmlExample1.html")
	if err != nil {
		return nil, err
	}
	data2, err := os.ReadFile("htmlfiles/htmlExample2.html")
	if err != nil {
		return nil, err
	}
	data3, err := os.ReadFile("htmlfiles/htmlExample3.html")
	if err != nil {
		return nil, err
	}

	dataFile1 := files.NewFile(string(data1))
	dataFile2 := files.NewFile(string(data2))
	dataFile3 := files.NewFile(string(data3))
	dataFiles := []*files.File{dataFile1, dataFile2, dataFile3}

	pk2 := ring.GeneratePrivate()
	u2 := user.NewUser(pk2)
	pk3 := ring.GeneratePrivate()
	u3 := user.NewUser(pk3)
	singers := []*user.Address{u2.Addr, u3.Addr}

	txs := make([]types.Transaction, 0)
	for _, d := range dataFiles {
		tx, err := transaction.NewTxPost(u, prevHash, d, singers)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

//Пример взятия постов из ноды
func ExampleGet(n *node) {
	fs, err := n.Get()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Files")
	for _, f := range fs {
		fmt.Println(f)
		fmt.Printf("\n\n\n")
	}
}

//Пример записи поста в ноду
func ExampleSet(n *node, post string) {
	if err := n.Set(post); err != nil {
		log.Fatal(err)
	}
}

func main() {
	n, err := NewNode()
	if err != nil {
		log.Fatal(err)
	}

	post := ".................some................"
	ExampleSet(n, post)
	ExampleGet(n)
	return
}
