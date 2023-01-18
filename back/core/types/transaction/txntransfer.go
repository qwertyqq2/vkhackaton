package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/qwertyqq2/filebc/user"

	"github.com/qwertyqq2/filebc/crypto"
)

type TxnTransfer struct {
	Rand      []byte `json:"rand"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Value     uint64 `json:"value"`
	ToStorage uint64 `json:"toStorage"`
	HashTx    []byte `json:"hashTx"`
	SignTx    []byte `json:"sign"`
	PrevBlock []byte `json:"prevBlock"`
}

func NewTxTransfer(sender *user.User, prevHash []byte, receiver *user.Address, value uint64) (*TxnTransfer, error) {
	rand := crypto.GenerateRandom()
	toStorage := uint64(1 * value / 10)
	tx := &TxnTransfer{
		Rand:      rand,
		Sender:    sender.Public(),
		Receiver:  receiver.String(),
		Value:     value,
		ToStorage: toStorage,
		PrevBlock: prevHash,
	}
	err := tx.Sign(sender)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TxnTransfer) hash() []byte {
	return crypto.HashSum(
		bytes.Join(
			[][]byte{
				t.Rand,
				crypto.Base64DecodeString(t.Sender),
				crypto.Base64DecodeString(t.Receiver),
				crypto.ToBytes(t.Value),
				crypto.ToBytes(t.ToStorage),
				t.PrevBlock,
			},
			[]byte{},
		))
}

func (t *TxnTransfer) Sign(u *user.User) error {
	h := t.hash()
	if h == nil {
		return fmt.Errorf("nil hash")
	}
	sign, err := u.SignData(h)
	if err != nil {
		return err
	}
	t.HashTx = h
	t.SignTx = sign
	return nil
}

func (t *TxnTransfer) hashValid() bool {
	return bytes.Equal(t.HashTx, t.hash())
}

func (t *TxnTransfer) signValid(senderstr string) bool {
	sender := user.ParseAddress(senderstr)
	if sender == nil {
		return false
	}
	err := user.VerifySign(sender, t.HashTx, t.SignTx)
	if err != nil {
		return false
	}
	return true
}

func (t *TxnTransfer) Valid() bool {
	if !t.hashValid() || !t.signValid(t.Sender) {
		return false
	}
	return true
}

func (t *TxnTransfer) SerializeTx() (string, error) {
	jsonData, err := json.MarshalIndent(*t, "", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
