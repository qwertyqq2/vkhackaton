package types

import (
	"bytes"
	"sort"
	"time"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

const (
	Difficulty = 15
)

type Block struct {
	PrevBlock values.Bytes      `json:"prevBlock"`
	PrevSnap  values.Bytes      `json:"prevSnap"`
	HashBlock values.Bytes      `json:"hashBlock"`
	CurShap   values.Bytes      `json:"curSnap"`
	Balances  map[string]uint64 `json:"balances"`
	Proof     uint64            `json:"nonce"`
	Time      string            `json:"time"`
	Miner     string            `json:"miner"`
	Diff      uint8             `json:"diff"`
	Sign      values.Bytes      `json:"sign"`

	transactions []Transaction
}

func NewBlock(prevBlock []byte, prevSnap []byte, miner *user.Address) *Block {
	return &Block{
		PrevBlock:    prevBlock,
		PrevSnap:     prevSnap,
		Miner:        miner.String(),
		Balances:     make(map[string]uint64),
		Diff:         Difficulty,
		transactions: make([]Transaction, 0),
	}
}

func NewGenesisBLock(creator *user.Address) *Block {
	gen := &Block{
		PrevBlock: []byte("GenBlock"),
		PrevSnap:  []byte("GenSnap"),
		Miner:     creator.String(),
		Balances:  make(map[string]uint64),
		Time:      time.Now().Format(time.RFC3339),
	}
	gen.Balances[creator.String()] = 100
	gen.HashBlock = gen.hash()
	gen.CurShap = gen.PrevSnap
	return gen
}

func (b *Block) hash() values.Bytes {
	temp := []byte{}
	for _, tx := range b.transactions {
		if tx.Hash() == nil {
			return nil
		}
		temp = values.HashSum(temp, tx.Hash())
	}
	list := []string{}
	for addr := range b.Balances {
		list = append(list, addr)
	}
	sort.Strings(list)
	for _, addr := range list {
		temp = values.HashSum(temp, []byte(addr), crypto.ToBytes(b.Balances[addr]))
	}
	return values.HashSum(temp,
		crypto.ToBytes(uint64(b.Diff)),
		b.PrevBlock,
		b.PrevSnap,
		[]byte(b.Miner),
		[]byte(b.Time))
}

func (b *Block) verifyHash() bool {
	return bytes.Equal(b.HashBlock, b.hash())
}

func (b *Block) verifySnapPrev() (bool, error) {
	coll, err := files.NewCollector()
	if err != nil {
		return false, err
	}
	snapPrev, err := coll.State()
	if err != nil {
		return false, err
	}
	return bytes.Equal(snapPrev, b.PrevSnap), nil
}
