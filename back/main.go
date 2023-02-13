package main

import (
	"fmt"

	"github.com/qwertyqq2/filebc/node"
)

func main() {
	n := node.NewNode()
	db := n.Get()
	for _, f := range db {
		fmt.Println(string(f))
	}
}
