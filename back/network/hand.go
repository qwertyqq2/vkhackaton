package network

import (
	"bufio"
	"log"

	"github.com/libp2p/go-libp2p/core/network"
)

type Conns map[string]Conn

type Handler struct {
	conns Conns
}

func NewHandler(conns Conns) Handler {
	return Handler{
		conns: conns,
	}
}

func (h Handler) run(s network.Stream) {
	h.handler(true)(s)
}

func (h Handler) handler(pend bool) func(s network.Stream) {
	return func(s network.Stream) {
		defer delete(h.conns, s.ID())
		h.conns[s.ID()] = NewConn(pend)
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go h.read(rw, s.ID())
		go h.write(rw, s.ID())
	}
}

func (h Handler) read(rw *bufio.ReadWriter, streamId string) {
	buf := make([]byte, 40*1024)
	for {
		n, err := rw.Read(buf)
		if err != nil {
			log.Println("Error reading from buffer")
			h.conns[streamId].In <- NilMessage()
			return
		}
		data := buf[:n]
		if data == nil {
			log.Println("nil data")
			return
		}
		msg, err := Unmarhsal(data)
		if err != nil {
			log.Println(err)
			return
		}
		if msg == nil {
			log.Println("nil msg")
			return
		}
		h.conns[streamId].In <- msg
	}
}

func (h Handler) write(rw *bufio.ReadWriter, streamId string) {
	for {
		select {
		case m := <-h.conns[streamId].Out:
			mar, err := Marhal(m)
			if err != nil {
				return
			}
			if _, err := rw.Write(mar); err != nil {
				return
			}
			if err := rw.Flush(); err != nil {
				return
			}
		}
	}
}
