package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"

	"github.com/qwertyqq2/filebc/crypto"
)

type TxnTransfer struct {
	Type      uint         `json:"type"`
	Rand      values.Bytes `json:"rand"`
	Sender    string       `json:"sender"`
	Receiver  string       `json:"receiver"`
	Value     uint64       `json:"value"`
	ToStorage uint64       `json:"toStorage"`
	HashTx    values.Bytes `json:"hashTx"`
	SignTx    values.Bytes `json:"sign"`
	PrevBlock values.Bytes `json:"prevBlock"`
}

func NewTxTransfer(sender *user.User, prevHash values.Bytes, receiver *user.Address, value uint64) (*TxnTransfer, error) {
	rand := crypto.GenerateRandom()
	toStorage := uint64(1 * value / 10)
	tx := &TxnTransfer{
		Type:      TypeTransferTx,
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

func (t *TxnTransfer) hash() values.Bytes {
	return values.HashSum(crypto.ToBytes(uint64(t.Type)),
		t.Rand,
		crypto.Base64DecodeString(t.Sender),
		crypto.Base64DecodeString(t.Receiver),
		crypto.ToBytes(t.Value),
		crypto.ToBytes(t.ToStorage),
		t.PrevBlock)
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

func (t *TxnTransfer) Hash() values.Bytes {
	return t.HashTx
}

func (t *TxnTransfer) DataTx() values.Bytes {
	return nil
}

func (t *TxnTransfer) SerializeTx() (string, error) {
	jsonData, err := json.MarshalIndent(*t, "", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializeTransferTx(data string) (*TxnTransfer, error) {
	var tx TxnTransfer
	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
