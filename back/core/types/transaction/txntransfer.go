package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/qwertyqq2/filebc/crypto/ring"
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
	Hash      values.Bytes `json:"hashTx"`
	Sign      values.Bytes `json:"sign"`
	PrevBlock values.Bytes `json:"prevBlock"`
}

func NewTxTransfer(sender *user.User, prevHash values.Bytes, receiver *user.Address, value uint64) (*TxnTransfer, error) {
	rand := crypto.GenerateRandom()
	toStorage := uint64(0)
	tx := &TxnTransfer{
		Type:      TypeTransferTx,
		Rand:      rand,
		Sender:    sender.Addr.String(),
		Receiver:  receiver.String(),
		Value:     value,
		ToStorage: toStorage,
		PrevBlock: prevHash,
	}
	err := tx.SignTx(sender)
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

func (t *TxnTransfer) SignTx(u *user.User) error {
	h := t.hash()
	if h == nil {
		return fmt.Errorf("nil hash")
	}
	sign, err := u.SignData(h)
	if err != nil {
		return err
	}
	t.Hash = h
	smarsh, err := sign.Marshal()
	if err != nil {
		return err
	}
	t.Sign = smarsh
	return nil
}

func (t *TxnTransfer) hashValid() bool {
	return bytes.Equal(t.Hash, t.hash())
}

func (t *TxnTransfer) signValid(senderstr string) bool {
	sender, err := user.ParseAddress(senderstr)
	if err != nil {
		return false
	}
	rign, err := ring.UnmarshalSign(t.Sign)
	if err != nil {
		return false
	}
	return ring.VerifySign(t.Hash, rign, sender.Public())
}

func (t *TxnTransfer) GetData() values.Bytes {
	return nil
}

func (t *TxnTransfer) GetHash() values.Bytes {
	return t.Hash
}

func (t *TxnTransfer) GetSender() string {
	return t.Sender
}

func (t *TxnTransfer) GetSenders() []string {
	return nil
}

func (t *TxnTransfer) GetReceiver() string {
	return t.Receiver
}

func (t *TxnTransfer) GetValue() uint64 {
	return t.Value
}

func (t *TxnTransfer) Valid() bool {
	if !t.hashValid() || !t.signValid(t.Sender) {
		return false
	}
	return true
}

func (t *TxnTransfer) Data() values.Bytes {
	return nil
}

func (t *TxnTransfer) Empty() error {
	if t.Sender == "" {
		return fmt.Errorf("nil sender")
	}
	if t.Receiver == "" {
		return fmt.Errorf("nil receiver")
	}
	if t.Value == 0 {
		return fmt.Errorf("nil value")
	}
	if t.Sign == nil {
		return fmt.Errorf("nil sign")
	}
	// if t.ToStorage == 0 {
	// 	return fmt.Errorf("nil storage")
	// }
	if t.Hash == nil {
		return fmt.Errorf("nil hash ")
	}
	return nil
}

func (t *TxnTransfer) GetType() uint {
	return t.Type
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
