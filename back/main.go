package main

import (
	"context"
	"flag"
	"log"

	nt "github.com/qwertyqq2/filebc/network"
)

func makeNode(port uint16) nt.P2PNode {
	conf := nt.DefaultConfig(port)

	node := nt.NewNode(*conf)

	ctx := context.Background()

	err := node.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return node
}

func main() {

	sourcePort := flag.Int("p", 0, "Source port number")
	list := flag.Bool("list", true, "listen or not")

	flag.Parse()

	node := makeNode(uint16(*sourcePort))
	if *list {
		err := node.Listen()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		err := node.Broadcast()
		if err != nil {
			log.Fatal(err)
		}
	}
	select {}
}
