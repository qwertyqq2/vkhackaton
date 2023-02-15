package node

import (
	"log"
	"os"

	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/network"
	"github.com/qwertyqq2/filebc/network/repo"
	"github.com/qwertyqq2/filebc/values"
)

type Option func(n string) error

type confNode struct {
	port       uint64
	repobcpath string
	confP2P    *network.ConfigNode
}

type Node struct {
	p2p network.P2PNode

	repo repo.Repo

	bc *core.Blockchain

	db []values.Bytes
}

func (n *Node) Init() error {
	return nil
}

func NewNode() *Node {
	data1, err := os.ReadFile("htmlfiles/htmlExample1.html")
	if err != nil {
		log.Fatal(err)
	}
	data2, err := os.ReadFile("htmlfiles/htmlExample2.html")
	if err != nil {
		log.Fatal(err)
	}
	data3, err := os.ReadFile("htmlfiles/htmlExample3.html")
	if err != nil {
		log.Fatal(err)
	}
	return &Node{
		db: []values.Bytes{data1, data2, data3},
	}
}

func (n *Node) Get() []values.Bytes {
	return n.db
}

// func NewNode(conf *confNode, opts ...Option) *Node {
// 	confP2P := conf.confP2P
// 	p2p := network.NewNode(*confP2P)
// 	ctx := context.Background()
// 	repo, err := reponode.Open("node-pk" + strconv.Itoa(int(conf.port)))

// 	var (
// 		wg      sync.WaitGroup
// 		errChan = make(chan error)
// 	)
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := p2p.Init(ctx); err != nil {
// 			log.Println(err)
// 			errChan <- err
// 		}
// 	}()

// 	pk, err := repo.PrivateKey()
// 	if err != nil {
// 		return nil
// 	}
// 	_ = user.NewUser(pk)
// 	return nil
// }
