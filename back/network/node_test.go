package network

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeNode(port uint16) P2PNode {
	conf := DefaultConfig(port)

	node := NewNode(*conf)

	ctx := context.Background()

	err := node.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return node
}

func TestNode(t *testing.T) {
	node := makeNode(3000)
	fmt.Println("Ok", node)
}

func TestMarhal(t *testing.T) {
	id := MsgName2
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
