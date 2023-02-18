package network

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peerstore"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
	"github.com/qwertyqq2/filebc/network/repo"
)

var (
	ErrNilBoostrap = errors.New("nil boostrap connect")
)

const (
	EnoughPeers = 1
)

type Conn struct {
	In            chan *Message
	Out           chan *Message
	Pending, Wait bool
}

func NewConn(pend bool) Conn {
	return Conn{
		In:      make(chan *Message),
		Out:     make(chan *Message),
		Pending: pend,
	}
}

type ConfigNode struct {
	Repopath          string
	Port              uint16
	Rendezvous        string
	ProtocolID        string
	LimitedConfigPath string
	BoostrapAddrs     []multiaddr.Multiaddr
}

func DefaultConfig(port uint16) *ConfigNode {
	return &ConfigNode{
		Repopath:          "repo-conf",
		Port:              port,
		Rendezvous:        "fbc",
		ProtocolID:        "/fbc/1.0.0",
		LimitedConfigPath: "limited-conf.json",
	}
}

type node struct {
	ConfigNode

	host host.Host
	repo repo.Repo

	addrs        []string
	boostrapInfo []peer.AddrInfo
	conns        Conns
	kadDHT       *dht.IpfsDHT

	streams map[peer.ID]bool

	sync.Mutex
	sync.WaitGroup
}

func NewNode(conf ConfigNode, conns map[string]Conn) P2PNode {
	return &node{
		host:       nil,
		ConfigNode: conf,
		conns:      conns,
		streams:    make(map[peer.ID]bool),
	}
}

func (n *node) ID() peer.ID {
	if n.host == nil {
		return ""
	}
	return n.host.ID()
}

func (n *node) Init(ctx context.Context) error {
	nodeAddrStrings := []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", n.Port)}
	repo, err := repo.Open("node-pk" + strconv.Itoa(int(n.Port)))
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
	limiterCfg, err := os.Open(n.LimitedConfigPath)
	if err != nil {
		return err
	}
	limiter, err := rcmgr.NewDefaultLimiterFromJSON(limiterCfg)
	if err != nil {
		return err
	}
	rcm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return err
	}
	it := 0
	var boostrapInfo []peer.AddrInfo
	if n.BoostrapAddrs == nil {
		boostrapInfo = dht.GetDefaultBootstrapPeerAddrInfos()
	} else {
		boostrapInfo = make([]peer.AddrInfo, len(n.BoostrapAddrs))
		for i := 0; i < len(boostrapInfo); i++ {
			inf, err := peer.AddrInfoFromP2pAddr(n.BoostrapAddrs[i])
			if err != nil {
				log.Println(err)
				continue
			}
			it++
			boostrapInfo[i] = *inf
		}
		if it == 0 {
			return fmt.Errorf("Nil info about addresses")
		}
	}
	n.boostrapInfo = boostrapInfo
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(nodeAddrStrings...),
		libp2p.EnableAutoRelay(
			autorelay.WithStaticRelays(boostrapInfo),
			autorelay.WithCircuitV1Support(),
		),
		libp2p.Identity(priv),
		libp2p.NATPortMap(),
		libp2p.ResourceManager(rcm),
	)
	if err != nil {
		return errors.Errorf("creating libp2p host error")
	}
	log.Println("Host created. We are:", host.Addrs()[0].String(), host.ID().Pretty())
	n.host = host
	p2pAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", host.ID().Pretty()))
	if err != nil {
		return errors.Errorf("creating host p2p multiaddr error")
	}
	host.SetStreamHandler(protocol.ID(n.ProtocolID), NewHandler(n.conns).handler(false))
	kademliaDHT, err := dht.New(
		ctx,
		host, dht.Mode(dht.ModeServer),
		//dht.RoutingTableRefreshPeriod(10*time.Second),
	)
	if err != nil {
		return err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	n.kadDHT = kademliaDHT
	var fullAddrs []string
	for _, addr := range host.Addrs() {
		fullAddrs = append(fullAddrs, addr.Encapsulate(p2pAddr).String())
	}
	n.addrs = fullAddrs
	return nil
}

// добавить освобождение памяти на релеере
func (n *node) boostrap(ctx context.Context, peerFindChan chan peer.AddrInfo, goClose <-chan bool, listener bool) error {
	var (
		it           = 0
		tctx, cancel = context.WithTimeout(ctx, time.Second*120)
		findPeer     = []peer.AddrInfo{}
	)
	defer cancel()
	for _, pinf := range n.boostrapInfo {
		n.host.Peerstore().AddAddrs(pinf.ID, pinf.Addrs, peerstore.PermanentAddrTTL)
		n.Add(1)
		go func(inf peer.AddrInfo) {
			if err := n.host.Connect(ctx, inf); err != nil {
				log.Println("Error connect to relay")
			} else {
				if err != nil {
					log.Printf("host failed to receive a relay reservation from relay. %v", err)
				} else {
					// _, err = client.Reserve(context.Background(), n.host, inf)
					// if err != nil {
					// 	log.Printf("host failed to receive a relay reservation from relay. %v", err)
					// } else {
					log.Println("Connection established with relay node")
					it++
					// }
				}
			}
			n.Done()
		}(pinf)
	}
	n.Wait()
	if it == 0 {
		return fmt.Errorf("Nil connect relay")
	}
	c, err := cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte("f"))
	if err != nil {
		return err
	}
	log.Println("Provider loading...")
	if err := n.kadDHT.Provide(tctx, c, true); err != nil {
		return fmt.Errorf("provide error")
	}
	log.Println("Provider declareted")
	if _, err := n.kadDHT.FindProviders(tctx, c); err != nil {
		return fmt.Errorf("find providers error")
	}
	routingDiscovery := drouting.NewRoutingDiscovery(n.kadDHT)
	dutil.Advertise(
		ctx,
		routingDiscovery,
		n.Rendezvous,
	)
	log.Println("Node search...")
	for {
		peerChan, err := routingDiscovery.FindPeers(
			ctx,
			n.Rendezvous,
			discovery.Limit(10),
			discovery.TTL(10*time.Microsecond),
		)
		if err != nil {
			return fmt.Errorf("find peers error")
		}
		exist := false
		for peer := range peerChan {
			if peer.ID == n.host.ID() {
				continue
			}
			for _, p := range findPeer {
				if p.ID == peer.ID {
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			peerFindChan <- peer
		}

		select {
		case <-goClose:
			return nil
		default:
		}
	}
}

func (n *node) Broadcast() error {
	var (
		it    = 0
		close = false

		peerFindChan = make(chan peer.AddrInfo, 10)
		goClose      = make(chan bool)
	)
	defer n.kadDHT.RoutingTable().Close()
	log.Println("Boostrap go")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(600*time.Second))
	defer cancel()
	go n.boostrap(context.Background(), peerFindChan, goClose, false)
	for !close {
		select {
		case <-ctx.Done():
			close = true
		case peer := <-peerFindChan:
			if _, ok := n.streams[peer.ID]; ok {
				break
			}
			n.streams[peer.ID] = true
			go func() {
				defer delete(n.streams, peer.ID)
				//log.Println("Connecting to:", peer.ID)
				stream, err := n.host.NewStream(network.WithUseTransient(context.Background(),
					n.ProtocolID), peer.ID, protocol.ID(n.ProtocolID))
				if err != nil {
					//log.Println("err conn: ", err)
					return
				} else {
					log.Println("connection established with anouther peer!")
					it++
					log.Println("Stream here")
					time.Sleep(1 * time.Second)
					NewHandler(n.conns).run(stream)
				}
			}()
			if it >= EnoughPeers {
				goClose <- true
				close = true
			}
		}
	}
	return nil
}

func (n *node) Listen() error {
	var (
		peerFindChan = make(chan peer.AddrInfo, 10)
		goClose      = make(chan bool)
	)
	go n.boostrap(context.Background(), peerFindChan, goClose, true)
	ctx, cancel := context.WithTimeout(context.Background(), 7200*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-peerFindChan:
			time.Sleep(6 * time.Second)
		}
	}
}

func (n *node) Addr() []string {
	return n.addrs
}