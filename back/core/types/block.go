package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

const (
	Difficulty = 10
)

type Block struct {
	Number    uint64       `json:"number"`
	PrevBlock values.Bytes `json:"prevBlock"`
	PrevSnap  values.Bytes `json:"prevSnap"`
	CurShap   values.Bytes `json:"curSnap"`
	HashBlock values.Bytes `json:"hashBlock"`
	Proof     values.Bytes `json:"proof"`
	Time      string       `json:"time"`
	Miner     string       `json:"miner"`
	Diff      uint8        `json:"diff"`
	Sign      values.Bytes `json:"sign"`
	TxsHash   values.Bytes `json:"txsHash"`

	accepted     bool
	transactions []Transaction
}

type Blocks []*Block

func NewBlock(prevNumber uint64, prevBlock, prevSnap values.Bytes, miner *user.Address, txs ...Transaction) *Block {
	return &Block{
		Number:       prevNumber + 1,
		PrevBlock:    prevBlock,
		PrevSnap:     prevSnap,
		Miner:        miner.String(),
		Diff:         Difficulty,
		transactions: txs,
		accepted:     false,
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
	return gen
}

func (b *Block) txhash() values.Bytes {
	temp := []byte{}
	for _, tx := range b.transactions {
		if tx.GetHash() == nil {
			return nil
		}
		temp = values.HashSum(temp, tx.GetHash())
	}
	return temp
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
		b.TxsHash,
		crypto.ToBytes(b.Number),
		b.PrevBlock,
		b.PrevSnap,
		b.CurShap,
		crypto.ToBytes(uint64(b.Diff)),
		[]byte(b.Miner),
		[]byte(b.Time))
}

func (b *Block) Accept(u *user.User) error {
	if b.accepted {
		return fmt.Errorf("block already accepted")
	}
	b.TxsHash = b.txhash()
	b.HashBlock = b.hash()
	s, err := u.SignData(b.HashBlock)
	if err != nil {
		return err
	}
	smar, err := s.Marshal()
	if err != nil {
		return err
	}
	b.Sign = smar
	proof, f := crypto.ProowOfWork(b.HashBlock, b.Diff, nil)
	if !f {
		return fmt.Errorf("Cant getting proof")
	}
	b.Proof = proof
	b.accepted = true
	return nil
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
	if block.CurShap == nil || block.HashBlock == nil || block.PrevBlock == nil {
		return fmt.Errorf("nil hash")
	}
	if block.Sign == nil {
		return fmt.Errorf("nil sign")
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

func (block *Block) TransactionsHash() values.Bytes {
	return block.TxsHash
}

func (block *Block) Accepted() bool {
	return block.accepted
}

func (block *Block) AllFiles() ([]*files.File, error) {
	fs := make([]*files.File, 0)
	for _, tx := range block.transactions {
		switch tx.GetType() {
		case transaction.TypePostTx:
			f, err := files.Deserialize(string(tx.GetData()))
			if err != nil {
				return nil, err
			}
			fs = append(fs, f)
		}
	}
	return fs, nil
}
