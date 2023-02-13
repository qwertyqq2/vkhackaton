package node

import (
	"log"
	"os"

	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/network"
	"github.com/qwertyqq2/filebc/values"
)

type Node struct {
	p2p network.P2PNode

	bc *core.Blockchain
	db []values.Bytes
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
