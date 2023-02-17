package node

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qwertyqq2/filebc/core"
	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/network"
	"github.com/qwertyqq2/filebc/user"
)

type Handler interface {
	GetChain() error

	GetBlocks() error

	SendChain() error

	GetTx() error

	SendTx() error

	SendBlocks() error
}

type Conn = network.Conn
type Conns = network.Conns

const (
	GetChain = iota
	GetBlocks
	SendChain
	SendBlocks
	GetTx
	SendTx
)

type handler struct {
	conns Conns

	client *user.Address
	bc     *core.Blockchain
}

func NewHandler(n *Node) handler {
	return handler{
		conns:  make(network.Conns),
		bc:     n.bc,
		client: n.client,
	}
}

func (h *handler) send(conn Conn, msgId int, payload []byte) {
	msg := network.NewMessage(msgId, payload)
	conn.In <- msg
}

func (h *handler) Send(msgId int, payload []byte) error {
	var (
		timeout = 300 * time.Second
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout sending")
		}
		for _, conn := range h.conns {
			if conn.Pending {
				h.send(conn, msgId, payload)
				go h.hand(conn)
				return nil
			}
		}
	}
	return nil
}

func (h *handler) listen() error {
	for {
		for _, conn := range h.conns {
			if !conn.Pending && !conn.Wait {
				go h.hand(conn)
				conn.Wait = true
			}
		}
	}
}

func (h *handler) hand(c Conn) error {
	for {
		select {

		case msg := <-c.Out:
			if network.IsNilMessage(msg) {
				return nil
			}
			switch msg.Id {
			case GetChain:
				d, err := h.SendChain()
				if err != nil {
					log.Println(err)
					return err
				}
				msg := network.NewMessage(SendChain, d)
				c.In <- msg

			case SendChain:
				bs, err := Decode(msg.Payload())
				if err != nil {
					log.Println(err)
				}
				blocks := make(types.Blocks, len(bs))
				for _, bss := range bs {
					b, err := types.DeserializeBlock(bss)
					if err != nil {
						log.Println(err)
						return err
					}
					blocks = append(blocks, b)
				}
				if _, err := core.NewBlockchainExternal(h.client, blocks...); err != nil {
					log.Println(err)
					return err
				}

			default:
				return fmt.Errorf("undefined msg")

				//
			}

		}
	}

}

func (h *handler) GetChain() ([]byte, error) {
	bs, err := h.bc.ReadChain()
	if err != nil {
		return nil, err
	}
	d, err := Encode(bs...)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (h *handler) SendChain() ([]byte, error) {
	bs, err := h.bc.ReadChain()
	if err != nil {
		return nil, err
	}
	d, err := Encode(bs...)
	if err != nil {
		return nil, err
	}
	return d, nil
}
