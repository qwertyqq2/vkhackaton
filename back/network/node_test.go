package network

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	node := makeNode(3000)
	fmt.Println("Ok", node)
}

func TestMarhal(t *testing.T) {
	id := MsgName3
	data := []byte("qeqweqw")
	msg := NewMessage(id, data)
	mar, err := Marhal(msg)
	if err != nil {
		t.Fatal("err marhal")
	}
	msgcopy, err := Unmarhsal(mar)
	if err != nil {
		t.Fatal("err unmarhal")
	}
	assert.Equal(t, msgcopy.Id, msg.Id)
	assert.Equal(t, msgcopy.payload, msg.payload)

}

func makeNode(port uint16) P2PNode {
	conf := DefaultConfig(port)

	node := NewNode(*conf, nil)

	ctx := context.Background()

	err := node.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return node
}
func TestNodeListen(t *testing.T) {

	sourcePort := 5000
	list := true

	flag.Parse()

	node := makeNode(uint16(sourcePort))
	if list {
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
