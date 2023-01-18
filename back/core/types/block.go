package types

import (
	"bytes"
	"sort"
	"time"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/user"
)

const (
	Difficulty = 15
)

type Block struct {
	PrevBlock []byte            `json:"prevBlock"`
	PrevSnap  []byte            `json:"prevSnap"`
	HashBlock []byte            `json:"hashBlock"`
	CurShap   []byte            `json:"curSnap"`
	Balances  map[string]uint64 `json:"balances"`
	Nonce     uint64            `json:"nonce"`
	Time      string            `json:"time"`
	Miner     string            `json:"miner"`
	Diff      uint8             `json:"diff"`
	Sign      []byte            `json:"sign"`

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

func (b *Block) hash() []byte {
	temp := []byte{}
	for _, tx := range b.transactions {
		if tx.Hash() == nil {
			return nil
		}
		temp = crypto.HashSum(
			bytes.Join(
				[][]byte{
					temp,
					tx.Hash(),
				},
				[]byte{},
			))
	}
	list := []string{}
	for addr := range b.Balances {
		list = append(list, addr)
	}
	sort.Strings(list)
	for _, addr := range list {
		temp = crypto.HashSum(
			bytes.Join(
				[][]byte{
					temp,
					[]byte(addr),
					crypto.ToBytes(b.Balances[addr]),
				},
				[]byte{},
			))
	}
	return crypto.HashSum(
		bytes.Join(
			[][]byte{
				temp,
				crypto.ToBytes(uint64(b.Diff)),
				b.PrevBlock,
				b.PrevSnap,
				[]byte(b.Miner),
				[]byte(b.Time),
			},
			[]byte{},
		))
}

func (b *Block) verifyHash() bool {
	return bytes.Equal(b.HashBlock, b.hash())
}
