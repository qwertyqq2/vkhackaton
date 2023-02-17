package node

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/network"
	reponode "github.com/qwertyqq2/filebc/node/repo"
	"github.com/qwertyqq2/filebc/user"
)

type Option func(n string) error

type confNode struct {
	port       uint16
	repobcpath string
	confP2P    *network.ConfigNode
}

func DefaultConf(port uint64) *confNode {
	return &confNode{
		port:       uint16(port),
		repobcpath: "repo",
		confP2P:    network.DefaultConfig(uint16(port)),
	}
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
	repo, err := reponode.Open("node-fbc-pk" + strconv.Itoa(int(conf.port)))
	if err != nil {
		log.Println(err)
		return nil
	}
	pk, err := repo.PrivateKey()
	if err != nil {
		return nil
	}
	log.Println("Your pk: ", crypto.Base64EncodeString(pk.Marshal()))
	client := user.NewUser(pk)
	n := &Node{
		repo: repo,
	}
	n.client = client.Address()
	n.hand = NewHandler(n)

	p2p := network.NewNode(*confP2P, n.hand.conns)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	if err := p2p.Init(ctx); err != nil {
		log.Println(err)
		return nil
	}
	n.p2p = p2p

	return n
}

func (n *Node) Listen() error {
	go n.p2p.Listen()
	go n.hand.listen()
	time.Sleep(1 * time.Second)
	return nil
}

func (n *Node) SendText() error {
	return n.Send(GetText, []byte("some"))
}

func (n *Node) Send(msgid int, data []byte) error {
	if err := n.p2p.Broadcast(); err != nil {
		log.Println("failed to query blockchain")
		return err
	}
	time.Sleep(1 * time.Second)
	if err := n.hand.send(msgid, data); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
