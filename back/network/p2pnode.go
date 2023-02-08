package network

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

type P2PNode interface {
	Init(context.Context) error

	ID() peer.ID

	Addr() []string

	Broadcast() error

	Listen() error
}
