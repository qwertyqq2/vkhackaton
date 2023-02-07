package network

import "github.com/libp2p/go-libp2p/core/peer"

type P2PNode interface {
	ID() peer.ID

	Broadcast() error
}
