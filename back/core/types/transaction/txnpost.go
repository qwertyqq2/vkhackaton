package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/crypto/ring"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/options"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

type TxnPost struct {
	Type      uint         `json:"type"`
	Rand      values.Bytes `json:"rand"`
	Senders   []string     `json:"senders"`
	PostId    values.Bytes `json:"postId"`
	ToStorage uint64       `json:"toStorage"`
	Hash      values.Bytes `json:"hashTx"`
	Seed      values.Bytes `json:"seed"`
	Signs     [][]byte     `json:"signs"`
	PrevBlock values.Bytes `json:"prevBlock"`
	Data      values.Bytes `json:"data"`
}

func NewTxPost(sender *user.User, prevHash values.Bytes, post *files.File, addrSigners []*user.Address) (*TxnPost, error) {
	rand := crypto.GenerateRandom()
	toStorage := post.Diff(options.MaxsizeFile)
	if !post.Verify(options.MaxsizeFile) {
		return nil, ErrIncorrectPost
	}
	tx := &TxnPost{
		Type:      TypePostTx,
		Rand:      rand,
		PrevBlock: prevHash,
		PostId:    post.Id,
		ToStorage: uint64(toStorage),
		Data:      post.Data,
	}
	pubsString := make([]string, len(addrSigners))
	for i := 0; i < len(pubsString); i++ {
		pubsString[i] = addrSigners[i].String()
	}
	tx.Senders = pubsString
	err := tx.SignTx(sender)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TxnPost) hash(addrs []string) values.Bytes {
	addr := make([][]byte, 0)
	for i := 0; i < len(addrs); i++ {
		addr = append(addr, crypto.Base64DecodeString(addrs[i]))
	}
	return values.HashSum(crypto.ToBytes(uint64(t.Type)),
		t.Rand,
		bytes.Join(addr, []byte{}),
		t.PrevBlock,
		t.PostId,
		crypto.ToBytes(t.ToStorage))
}

func (t *TxnPost) SignTx(u *user.User) error {
	round := RandomIntRange(1, len(t.Senders)-1)
	t.Senders = append(t.Senders[:round+1], t.Senders[:round]...)
	t.Senders[round] = u.Addr.String()
	pubs := make([]*ring.PublicKey, len(t.Senders))
	h := t.hash(t.Senders)
	if h == nil {
		return fmt.Errorf("nil hash")
	}
	for i := 0; i < len(pubs); i++ {
		addr, err := user.ParseAddress(t.Senders[i])
		if err != nil {
			return err
		}
		pubs[i] = addr.Public()
	}
	ringSign, err := u.RingSignData(h, pubs, round)
	if err != nil {
		return err
	}
	t.Hash = h
	t.Seed = ringSign.Seed
	t.Signs = ringSign.Sings
	return nil
}

func (t *TxnPost) hashValid() bool {
	return bytes.Equal(t.Hash, t.hash(t.Senders))
}

func (t *TxnPost) signValid() bool {
	senders := make([]*user.Address, len(t.Senders))
	for i := 0; i < len(senders); i++ {
		s, err := user.ParseAddress(t.Senders[i])
		if err != nil {
			return false
		}
		senders[i] = s
	}
	return user.VeryfySignRing(t.Hash, senders, t.Seed, t.Signs)
}

func (t *TxnPost) GetData() values.Bytes {
	return t.Data
}

func (t *TxnPost) GetHash() values.Bytes {
	return t.Hash
}

func (t *TxnPost) GetSenders() []string {
	return t.Senders
}

func (t *TxnPost) GetReceiver() string {
	return ""
}

func (t *TxnPost) GetValue() uint64 {
	return 0
}

func (t *TxnPost) GetSender() string {
	return ""
}

func (t *TxnPost) Valid() bool {
	if !t.hashValid() || !t.signValid() {
		return false
	}
	return true
}

func (t *TxnPost) Empty() error {
	if len(t.Senders) == 0 {
		return fmt.Errorf("nil len senders")
	}
	for _, s := range t.Senders {
		if s == "" {
			return fmt.Errorf("nil senders")
		}
	}
	if t.Data == nil {
		return fmt.Errorf("nil data post")
	}
	if len(t.Signs) == 0 {
		return fmt.Errorf("nil len signs")
	}

	for _, s := range t.Signs {
		if s == nil {
			return fmt.Errorf("nil sign")
		}
	}
	if t.ToStorage == 0 {
		return fmt.Errorf("nil storage")
	}
	if t.Hash == nil {
		return fmt.Errorf("nil hash ")
	}
	return nil
}

func (t *TxnPost) GetType() uint {
	return t.Type
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

func RandomIntRange(m, n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n-m+1) + m
}
