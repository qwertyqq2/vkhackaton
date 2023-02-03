package listener

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/qwertyqq2/filebc/network/repo"
)

var (
	ErrNilBoostrap = errors.New("nil boostrap connect")
)

type P2PNode interface {
}

type ConfigNode struct {
	Repopath      string
	Port          uint16
	BoostrapAddrs []multiaddr.Multiaddr
}

type node struct {
	ConfigNode

	host host.Host
	repo repo.Repo

	addrs  []string
	loaded bool

	kadDHT *dht.IpfsDHT

	hanler func(*network.Stream)

	sync.Mutex
	wg sync.WaitGroup
}

func NewNode(conf ConfigNode) *node {
	return &node{
		host:       nil,
		ConfigNode: conf,
	}
}

func (n *node) ID() peer.ID {
	if n.host == nil {
		return ""
	}
	return n.host.ID()
}

func (n *node) Init(ctx context.Context, boostrapOnly bool) error {
	nodeAddrStrings := []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", n.Port)}
	repo, err := repo.Open("node-pk" + string(n.Port))
	if err != nil {
		return err
	}
	priv, err := repo.PrivateKey()
	if err != nil {
		return err
	}
	n.repo = repo
	if err != nil {
		return err
	}
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(nodeAddrStrings...),
		libp2p.Identity(priv),
	)
	if err != nil {
		return errors.Errorf("creating libp2p host error")
	}
	n.host = host

	p2pAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", host.ID().Pretty()))
	if err != nil {
		return errors.Errorf("creating host p2p multiaddr error")
	}

	var fullAddrs []string
	for _, addr := range host.Addrs() {
		fullAddrs = append(fullAddrs, addr.Encapsulate(p2pAddr).String())
	}
	n.addrs = fullAddrs
	fmt.Println("Addresses : ", fullAddrs)
	if boostrapOnly {
		if len(n.BoostrapAddrs) == 0 {
			return fmt.Errorf("nil boostrap peers")
		}
		return n.Boostrap(ctx)
	}
	return nil
}

func (n *node) Boostrap(ctx context.Context) error {
	var boostrapPeers []peer.AddrInfo
	for _, addr := range n.BoostrapAddrs {
		pinf, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		n.host.Peerstore().AddAddrs(pinf.ID, pinf.Addrs, peerstore.PermanentAddrTTL)
		boostrapPeers = append(boostrapPeers, *pinf)
	}
	kadDHT, err := dht.New(
		ctx,
		n.host,
		dht.BootstrapPeers(boostrapPeers...),
		dht.ProtocolPrefix("pref"),
		dht.Mode(dht.ModeAutoServer),
	)
	n.kadDHT = kadDHT
	if err != nil {
		return fmt.Errorf("new dht error")
	}
	log.Println("Boostrap begin")
	if err := kadDHT.Bootstrap(ctx); err != nil {
		return fmt.Errorf("boostrap dht error")
	}
	n.wg.Add(1)
	var (
		it = 0
	)
	log.Println("Connect to peers")
	for _, pinf := range boostrapPeers {
		go func(p peer.AddrInfo) {
			if err := n.host.Connect(ctx, p); err != nil {
				log.Println(err)
				return
			}
			it++
		}(pinf)
	}
	n.wg.Wait()
	if it == 0 {
		return ErrNilBoostrap
	}
	log.Println("Set handler")
	n.host.SetStreamHandler("some pid", func(s network.Stream) {
		fmt.Println("Stream has been received")
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go read(rw)
		go write(rw)
	})
	n.loaded = true
	return nil
}

func read(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func write(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func (n *node) RunStream(ctx context.Context, targetAddr multiaddr.Multiaddr) error {
	info, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		return err
	}
	fmt.Println("Start stream")
	_, err = n.host.NewStream(context.Background(), info.ID, "some pid")
	if err != nil {
		return err
	}
	return nil
}

func (n *node) Addr() []string {
	return n.addrs
}
