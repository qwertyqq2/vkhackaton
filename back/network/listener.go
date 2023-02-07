package network

import (
	"bufio"
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

type ConfigNode struct {
	Repopath          string
	Port              uint16
	Rendezvous        string
	ProtocolID        string
	LimitedConfigPath string
	BoostrapAddrs     []multiaddr.Multiaddr
}

type node struct {
	ConfigNode

	host host.Host
	repo repo.Repo

	addrs        []string
	loaded       bool
	boostrapInfo []peer.AddrInfo

	kadDHT *dht.IpfsDHT

	hanler func(*network.Stream)

	sync.Mutex
	sync.WaitGroup
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

func (n *node) Init(ctx context.Context) (P2PNode, error) {
	nodeAddrStrings := []string{fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", n.Port)}
	repo, err := repo.Open("node-pk" + strconv.Itoa(int(n.Port)))
	if err != nil {
		return nil, err
	}
	priv, err := repo.PrivateKey()
	if err != nil {
		return nil, err
	}
	n.repo = repo
	if err != nil {
		return nil, err
	}
	limiterCfg, err := os.Open(n.LimitedConfigPath)
	if err != nil {
		return nil, err
	}
	limiter, err := rcmgr.NewDefaultLimiterFromJSON(limiterCfg)
	if err != nil {
		return nil, err
	}
	rcm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}
	it := 0
	boostrapInfo := make([]peer.AddrInfo, len(n.BoostrapAddrs))
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
		return nil, fmt.Errorf("Nil info about addresses")
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
		return nil, errors.Errorf("creating libp2p host error")
	}
	n.host = host
	p2pAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", host.ID().Pretty()))
	if err != nil {
		return nil, errors.Errorf("creating host p2p multiaddr error")
	}
	kademliaDHT, err := dht.New(ctx, host, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var fullAddrs []string
	for _, addr := range host.Addrs() {
		fullAddrs = append(fullAddrs, addr.Encapsulate(p2pAddr).String())
	}
	n.addrs = fullAddrs
	return n, nil
}

// добавить освобождение памяти на релеере
func (n *node) boostrap(ctx context.Context, peerFindChan chan peer.AddrInfo, goClose <-chan bool) error {
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
					log.Println("Connection established with relay node")
					it++
				}
			}
			n.Done()
		}(pinf)
	}
	n.Wait()
	if it == 0 {
		return fmt.Errorf("Nil connect relay")
	}
	c, err := cid.NewPrefixV1(cid.Raw, multihash.SHA2_256).Sum([]byte("meet me here"))
	if err != nil {
		return err
	}
	if err := n.kadDHT.Provide(tctx, c, true); err != nil {
		return fmt.Errorf("provide error")
	}
	if _, err := n.kadDHT.FindProviders(tctx, c); err != nil {
		return fmt.Errorf("find providers error")
	}
	routingDiscovery := drouting.NewRoutingDiscovery(n.kadDHT)
	dutil.Advertise(ctx, routingDiscovery, n.Rendezvous)
	for {
		go func() {
			peerChan, err := routingDiscovery.FindPeers(
				ctx,
				n.Rendezvous,
				discovery.Limit(100),
			)
			if err != nil {
				return
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
		}()
		select {
		case <-goClose:
			return nil
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (n *node) Broadcast() error {
	var (
		wg    sync.WaitGroup
		it    = 0
		close = false

		peerFindChan = make(chan peer.AddrInfo)
		goClose      = make(chan bool)
	)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	go n.boostrap(ctx, peerFindChan, goClose)
	for !close {
		select {
		case <-ctx.Done():
			close = true
		case peer := <-peerFindChan:
			go func(it int) {
				it++
				stream, err := n.host.NewStream(network.WithUseTransient(context.Background(),
					n.ProtocolID), peer.ID, protocol.ID(n.ProtocolID))
				if err != nil {
					return
				} else {
					wg.Add(1)
					rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
					closed := make(chan bool)
					go writeData(rw, closed, "biba")
					go readData(rw, closed)
					select {
					case <-closed:
						wg.Done()
						return
					}
				}
			}(it)
			if it >= EnoughPeers {
				goClose <- true
			}
		}
	}
	wg.Wait()

	return nil
}

func (n *node) Addr() []string {
	return n.addrs
}
