package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/qwertyqq2/filebc/node"
)

var port uint

func main() {

	port := flag.Uint64("p", 4000, "your port ")
	list := flag.Bool("l", false, "listen or not")
	flag.Parse()

	node := node.NewNode(node.DefaultConf(*port))

	if *list {
		fmt.Println("listen")
		node.Listen()
	} else {
		if err := node.SendText(); err != nil {
			log.Fatal(err)
		}
	}
	select {}

}
