package node

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/network"
	reponode "github.com/qwertyqq2/filebc/node/repo"
	"github.com/qwertyqq2/filebc/user"
)

type Option func(n string) error

type confNode struct {
	port       uint64
	repobcpath string
	confP2P    *network.ConfigNode
}

type Node struct {
	p2p network.P2PNode

	hand handler

	client *user.Address

	repo reponode.Repo

	bc *core.Blockchain
}

func NewNode(conf *confNode, opts ...Option) *Node {
	confP2P := conf.confP2P
	repo, err := reponode.Open("node-pk" + strconv.Itoa(int(conf.port)))
	if err != nil {
		log.Println(err)
		return nil
	}
	pk, err := repo.PrivateKey()
	if err != nil {
		return nil
	}
	client := user.NewUser(pk)
	n := &Node{
		repo: repo,
	}
	n.client = client.Address()
	n.hand = NewHandler(n)

	p2p := network.NewNode(*confP2P, n.hand.conns)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// var (
	// 	wg      sync.WaitGroup
	// 	errChan = make(chan error)
	// )

	if err := p2p.Init(ctx); err != nil {
		log.Println(err)
		return nil
	}

	// bc, err := core.LoadBlockchain(client.Address())
	// if err != nil {
	// 	switch err {
	// 	case core.ErrLoadBc:
	// 		if err := p2p.Broadcast(); err != nil {
	// 			log.Println("failed to query blockchain")
	// 			return nil
	// 		}
	// 		go n.hand.Send(GetChain, []byte(""))

	// 	default:
	// 		log.Println("failed to load blockchain")
	// 		return nil
	// 	}
	// }
	if err := p2p.Broadcast(); err != nil {
		log.Println("failed to query blockchain")
		return nil
	}
	go n.hand.Send(GetChain, []byte(""))

	return nil
}
