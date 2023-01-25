package types

import (
	"encoding/json"
	"fmt"
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
	Proof     values.Bytes `json:"nonce"`
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
		if tx.GetHash() == nil {
			return nil
		}
		temp = values.HashSum(temp, tx.GetHash())
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
		data[i] = tx.GetData()
	}
	return data
}

func (block *Block) Pow() error {
	nonce, f := crypto.ProowOfWork(block.CurShap, block.Diff, nil)
	if !f {
		return fmt.Errorf("cant pow")
	}
	block.Proof = nonce
	return nil
}

func (block *Block) EmptyBlock() error {
	if block.CurShap == nil || block.HashBlock == nil || block.PrevBlock == nil || block.Sign == nil {
		return fmt.Errorf("nil hash")
	}
	if block.Number == 0 {
		return fmt.Errorf("zero number")
	}
	if block.Miner == "" {
		return fmt.Errorf("nil miner")
	}
	if len(block.transactions) == 0 {
		return fmt.Errorf("0 txs")
	}
	return nil
}

func (block *Block) SerializeBlock() (string, error) {
	jsonData, err := json.MarshalIndent(block, " ", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializeBlock(data string) (*Block, error) {
	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (block *Block) Transactions() []Transaction {
	return block.transactions
}

type blockIterator struct {
	index        int
	transactions []Transaction
}

func BlockIterator(block *Block) *blockIterator {
	return &blockIterator{
		index:        -1,
		transactions: block.transactions,
	}
}

func (iter *blockIterator) next() (Transaction, error) {
	if iter.index+1 >= len(iter.transactions) {
		iter.index = len(iter.transactions)
		return nil, fmt.Errorf("end iterator")
	}
	iter.index++
	return iter.transactions[iter.index], nil
}

func (iter *blockIterator) prev() Transaction {
	if iter.index < 1 {
		return nil
	}
	return iter.transactions[iter.index-1]
}

func (iter *blockIterator) current() Transaction {
	if iter.index == -1 || iter.index >= len(iter.transactions) {
		return nil
	}
	return iter.transactions[iter.index]
}

func (iter *blockIterator) first() Transaction {
	return iter.transactions[0]
}

func (iter *blockIterator) remaining() int {
	return len(iter.transactions) - iter.index
}

func (iter *blockIterator) processed() int {
	return iter.index + 1
}
