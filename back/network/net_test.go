package network

import (
	"fmt"
	"net"
	"testing"
)

func TestConn(t *testing.T) {
	addr := "127.0.0.1"
	port := "8888"
	listener, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ip: ", host)
	fmt.Println("port: ", port)
}
