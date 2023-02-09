package node

import (
	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/network"
)

type Node struct {
	p2p network.P2PNode

	bc *core.Blockchain
}
