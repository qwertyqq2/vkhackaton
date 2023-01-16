package network

import (
	"log"
	"net"
	"strings"
)

const (
	EndBytes = "\0005\0000\0001"
	Waitime  = 10
	BufSize  = 4 << 10
	MaxSize  = 2 << 20
)

type Listener net.Listener
type Conn net.Conn

type Node struct {
	addr string

	conn     map[Conn]bool
	listener Listener
	listen   bool
}

func NewNode(addr string) *Node {
	return &Node{
		addr:   addr,
		conn:   make(map[Conn]bool),
		listen: false,
	}
}

func (n *Node) Listen(port string, handle func(Conn, *Package)) (Listener, error) {
	listener, err := net.Listen("tcp", n.addr+":"+port)
	if err != nil {
		return nil, err
	}
	n.listen = true
	n.listener = listener
	go n.serve(handle)
	return Listener(listener), nil
}

func (n *Node) Close() error {
	err := n.listener.Close()
	if err != nil {
		return err
	}
	n.listen = false
	for c := range n.conn {
		delete(n.conn, c)
	}
	return nil
}

func (n *Node) serve(handle func(Conn, *Package)) {
	defer n.Close()
	for {
		c, err := n.listener.Accept()
		if err != nil {
			return
		}
		go func(conn net.Conn, handle func(Conn, *Package)) {
			n.conn[c] = true
			defer conn.Close()
			pack := readPack(conn)
			if pack == nil {
				return
			}
			handle(Conn(conn), pack)
			n.conn[c] = false
		}(c, handle)
	}
}

func readPack(conn net.Conn) *Package {
	var (
		size = uint64(0)
		buf  = make([]byte, BufSize)
		data string
	)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return nil
		}
		size += uint64(length)
		if size > MaxSize {
			log.Println(err)
			return nil
		}
		data += string(buf[:length])
		if strings.Contains(data, EndBytes) {
			data = strings.Split(data, EndBytes)[0]
			break
		}
	}
	deserializePack, err := DeserializePack(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return deserializePack
}
