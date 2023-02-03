package main

import (
	"context"
	"flag"
	"log"

	"github.com/multiformats/go-multiaddr"
	li "github.com/qwertyqq2/filebc/network/listener"
)

func main() {

	sourcePort := flag.Int("p", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")

	flag.Parse()

	if *dest == "" {
		n := li.NewNode(
			li.ConfigNode{
				Port: uint16(*sourcePort),
			},
		)
		err := n.Init(context.Background(), false)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatal(err)
		}
		n := li.NewNode(
			li.ConfigNode{
				Port:          uint16(*sourcePort),
				BoostrapAddrs: []multiaddr.Multiaddr{maddr},
			},
		)
		err = n.Init(context.Background(), true)
		if err != nil {
			log.Fatal(err)
		}
		if err := n.RunStream(context.Background(), maddr); err != nil {
			log.Fatal(err)
		}
	}

	select {}
}
