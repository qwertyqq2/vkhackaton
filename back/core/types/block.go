package types

import (
	"time"

	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

const (
	Difficulty = 15
)

type Block struct {
	Number    uint64       `json:"number"`
	PrevBlock values.Bytes `json:"prevBlock"`
	PrevSnap  values.Bytes `json:"prevSnap"`
	CurShap   values.Bytes `json:"curSnap"`
	HashBlock values.Bytes `json:"hashBlock"`
	Proof     uint64       `json:"nonce"`
	Time      string       `json:"time"`
	Miner     string       `json:"miner"`
	Diff      uint8        `json:"diff"`
	Sign      values.Bytes `json:"sign"`

	transactions []Transaction
}

type Blocks []*Block

func NewBlock(prevNumber uint64, prevBlock, prevSnap values.Bytes, miner *user.Address) *Block {
	return &Block{
		Number:       prevNumber + 1,
		PrevBlock:    prevBlock,
		PrevSnap:     prevSnap,
		Miner:        miner.String(),
		Diff:         Difficulty,
		transactions: make([]Transaction, 0),
	}
}

func NewGenesisBLock(creator *user.Address) *Block {
	gen := &Block{
		Number:    1,
		PrevBlock: []byte("GenBlock"),
		PrevSnap:  []byte("GenSnap"),
		Miner:     creator.String(),
		Time:      time.Now().Format(time.RFC3339),
	}
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
	return values.HashSum(
		temp,
		crypto.ToBytes(b.Number),
		b.PrevBlock,
		b.PrevSnap,
		b.CurShap,
		crypto.ToBytes(uint64(b.Diff)),
		[]byte(b.Miner),
		[]byte(b.Time))
}

func (block *Block) Data() []values.Bytes {
	data := make([]values.Bytes, len(block.transactions))
	for i, tx := range block.transactions {
		data[i] = tx.Data()
	}
	return data
}
