package core

import (
	"fmt"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/files"
	"github.com/qwertyqq2/filebc/user"
	"github.com/qwertyqq2/filebc/values"
)

var (
	MinPostSize     = 100
	MaxPostSize     = 10000
	MinTokenReserve = 50
	Fees            = 10
)

type validator struct {
	bc *Blockchain

	state values.Bytes
}

func newValidator(bc *Blockchain) *validator {
	return &validator{
		bc:    bc,
		state: bc.snap,
	}
}

func (validator *validator) add(txs ...types.Transaction) error {
	for _, tx := range txs {
		if err := tx.Empty(); err != nil {
			return err
		}
		switch tx.GetType() {
		case transaction.TypePostTx:
			if !validator.validPostSize(tx.GetData()) {
				return fmt.Errorf("invalid size post")
			}
			if !validator.validMinReserveForPost(tx.GetSender()) {
				return fmt.Errorf("not enough tokens for post")
			}
			file := files.NewFile(string(tx.GetData()))
			state, err := validator.bc.coll.State(file)
			if err != nil {
				return err
			}
			validator.state = state

		case transaction.TypeTransferTx:

		}
	}
	return nil
}

func (validator *validator) validPostSize(post values.Bytes) bool {
	return len(post) > MinPostSize && len(post) < MaxPostSize
}

func (validator *validator) validMinReserveForPost(sender string) bool {
	bal, err := validator.bc.coll.Balance(user.ParseAddress(sender))
	if err != nil {
		return false
	}
	return bal > uint64(MinTokenReserve)
}

func (validator *validator) validValue(sender string, value uint64) bool {
	bal, err := validator.bc.coll.Balance(user.ParseAddress(sender))
	if err != nil {
		return false
	}
	return bal > value+value*uint64(Fees)/100
}
