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
	GetText
	SendText
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

func (h *handler) listen() error {
	for {
		//time.Sleep(1 * time.Second)
		for _, conn := range h.conns {
			if conn.Wait {
				conn.Wait = false
				go h.hand(conn)
			}
		}
	}
}

func (h *handler) send(msgId int, payload []byte) error {
	var (
		timeout = 300 * time.Second
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout sending")
		default:
		}
		for _, conn := range h.conns {
			fmt.Println(conn.ID, conn.Pending, conn.Wait)
		}
		time.Sleep(1 * time.Second)
		for _, conn := range h.conns {
			if conn.Pending {
				conn.Pending = false
				msg := network.NewMessage(msgId, payload)
				conn.In <- msg
				go h.hand(conn)
				return nil
			}
		}
	}
}

func (h *handler) hand(c *Conn) error {
	for {
		select {
		case msg := <-c.Out:
			if network.IsNilMessage(msg) {
				log.Println("rec nil msg")
				time.Sleep(2 * time.Second)
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
				c.In <- network.NilMessage()
				return nil

			case GetText:
				msg := network.NewMessage(SendText, []byte("????????"))
				c.In <- msg

			case SendText:
				fmt.Println("Text here")
				c.In <- network.NilMessage()
				return nil

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
