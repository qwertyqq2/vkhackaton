package types

import (
	"encoding/json"
	"fmt"
	"os/user"

	"github.com/qwertyqq2/filebc/core/types/transaction"
)

type Transaction interface {
	hash() ([]byte, error)

	Sign(*user.User) []byte

	Valid() bool

	SerializeTx() (string, error)
}

func DeserializeTx(stx string) (interface{}, error) {
	var tx interface{}
	err := json.Unmarshal([]byte(stx), &tx)
	if err != nil {
		return nil, err
	}
	post, okpost := tx.(transaction.TxnPost)
	trans, oktrans := tx.(transaction.TxnTransfer)
	if okpost {
		fmt.Println("Ok post")
		return post, nil
	}
	if oktrans {
		fmt.Println("Ok trans")
		return trans, nil
	}
	fmt.Println("Not convert")
	return nil, nil
}
