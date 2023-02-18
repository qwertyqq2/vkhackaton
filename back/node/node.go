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
	Conf *confNode

	p2p network.P2PNode

	hand handler

	client *user.Address

	repo reponode.Repo

	bc *core.Blockchain
}

func NewNode(conf *confNode, opts ...Option) *Node {
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
		Conf: conf,
	}
	n.client = client.Address()
	n.hand = NewHandler(n)

	return n
}

func (n *Node) Listen() error {
	p2p := network.NewNode(*n.Conf.confP2P, n.hand.conns)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	if err := p2p.Init(ctx); err != nil {
		log.Println(err)
		return nil
	}
	go n.hand.listen()
	p2p.Listen()
	//time.Sleep(1 * time.Second)
	return nil
}

func (n *Node) SendText() error {
	return n.Send(GetText, []byte("some"))
}

func (n *Node) Send(msgid int, data []byte) error {
	p2p := network.NewNode(*n.Conf.confP2P, n.hand.conns)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	if err := p2p.Init(ctx); err != nil {
		log.Println(err)
		return nil
	}
	if err := p2p.Broadcast(); err != nil {
		log.Println("failed to query blockchain")
		return err
	}
	time.Sleep(2 * time.Second)
	log.Println("Sending")
	if err := n.hand.send(msgid, data); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
