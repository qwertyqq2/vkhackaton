package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/options"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

type TxnPost struct {
	Type      uint         `json:"type"`
	Rand      values.Bytes `json:"rand"`
	Sender    string       `json:"sender"`
	PostId    values.Bytes `json:"postId"`
	ToStorage uint64       `json:"toStorage"`
	HashTx    values.Bytes `json:"hashTx"`
	SignTx    values.Bytes `json:"signTx"`
	PrevBlock values.Bytes `json:"prevBlock"`
	DataTx    values.Bytes `json:"data"`
}

func NewTxPost(sender *user.User, prevHash values.Bytes, post *files.File) (*TxnPost, error) {
	rand := crypto.GenerateRandom()
	toStorage := post.Diff(options.MaxsizeFile)
	if !post.Verify(options.MaxsizeFile) {
		return nil, ErrIncorrectPost
	}
	tx := &TxnPost{
		Type:      TypePostTx,
		Rand:      rand,
		Sender:    sender.Public(),
		PrevBlock: prevHash,
		PostId:    post.Id,
		ToStorage: uint64(toStorage),
		DataTx:    post.Data,
	}
	err := tx.Sign(sender)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TxnPost) hash() values.Bytes {
	return values.HashSum(crypto.ToBytes(uint64(t.Type)),
		t.Rand,
		crypto.Base64DecodeString(t.Sender),
		t.PrevBlock,
		t.PostId,
		crypto.ToBytes(t.ToStorage))
}

func (t *TxnPost) Sign(u *user.User) error {
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

func (t *TxnPost) hashValid() bool {
	return bytes.Equal(t.HashTx, t.hash())
}

func (t *TxnPost) signValid(senderstr string) bool {
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

func (t *TxnPost) Valid() bool {
	if !t.hashValid() || !t.signValid(t.Sender) {
		return false
	}
	return true
}

func (t *TxnPost) Hash() values.Bytes {
	return t.HashTx
}

func (t *TxnPost) Data() values.Bytes {
	return t.DataTx
}

func (t *TxnPost) SerializeTx() (string, error) {
	jsonData, err := json.MarshalIndent(*t, "", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializePostTx(data string) (*TxnPost, error) {
	var tx TxnPost
	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
