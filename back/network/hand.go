package network

import (
	"bufio"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/network"
)

type Conns map[string]*Conn

type Handler struct {
	conns Conns
}

func NewHandler(conns Conns) Handler {
	return Handler{
		conns: conns,
	}
}

func (h Handler) run(s network.Stream) {
	h.handler(true, false)(s)
}

func (h Handler) handler(pend, wait bool) func(s network.Stream) {
	return func(s network.Stream) {
		log.Println("Стрим прищел")
		h.conns[s.ID()] = NewConn(pend, wait, s.ID())
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go h.read(rw, s.ID())
		go h.write(rw, s.ID())
	}
}

func (h Handler) read(rw *bufio.ReadWriter, streamId string) {
	defer func() {
		fmt.Println("Read err")
		delete(h.conns, streamId)
	}()
	buf := make([]byte, 40*1024)
	for {
		n, err := rw.Read(buf)
		if err != nil {
			log.Println("Error reading from buffer")
			h.conns[streamId].Out <- NilMessage()
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
		h.conns[streamId].Out <- msg
	}
}

func (h Handler) write(rw *bufio.ReadWriter, streamId string) {
	for {
		select {
		case m := <-h.conns[streamId].In:
			if IsNilMessage(m) {
				return
			}
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
