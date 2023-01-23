package types

import (
	"encoding/json"
	"fmt"

	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"

	"github.com/qwertyqq2/filebc/core/types/transaction"
)

type Transaction interface {
	SignTx(*user.User) error

	Valid() bool

	GetHash() values.Bytes

	GetSender() string

	GetReceiver() string

	GetValue() uint64

	GetData() values.Bytes

	GetType() uint

	Empty() error

	SerializeTx() (string, error)
}

func DeserializeTx(stx string) (Transaction, error) {
	var tx interface{}
	err := json.Unmarshal([]byte(stx), &tx)
	if err != nil {
		return nil, err
	}
	m := tx.(map[string]interface{})
	var t uint
	if tstr, ok := m["type"]; ok {
		t = uint(tstr.(float64))
	}
	switch t {
	case transaction.TypePostTx:
		var txp transaction.TxnPost
		err := json.Unmarshal([]byte(stx), &txp)
		if err != nil {
			return nil, err
		}
		return &txp, nil
	case transaction.TypeTransferTx:
		var txtrans transaction.TxnTransfer
		err := json.Unmarshal([]byte(stx), &txtrans)
		if err != nil {
			return nil, err
		}
		return &txtrans, nil
	}
	return nil, fmt.Errorf("undefined type tx")
}
