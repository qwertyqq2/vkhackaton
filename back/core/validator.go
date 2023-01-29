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

func (validator *validator) add(state values.Bytes, txs ...types.Transaction) (values.Bytes, error) {
	for _, tx := range txs {
		if err := tx.Empty(); err != nil {
			return nil, err
		}
		switch tx.GetType() {
		case transaction.TypePostTx:
			if !validator.validPostSize(tx.GetData()) {
				return nil, fmt.Errorf("invalid size post")
			}
			if !validator.validMinReserveForPost(tx.GetSender()) {
				return nil, fmt.Errorf("not enough tokens for post")
			}
			file := files.NewFile(string(tx.GetData()))
			return validator.bc.coll.AddFile(state, file), nil

		case transaction.TypeTransferTx:
			if !validator.validValue(tx.GetSender(), tx.GetValue()) {
				return nil, fmt.Errorf("invalid value")
			}
			sender, err := user.ParseAddress(tx.GetSender())
			if err != nil {
				return nil, fmt.Errorf("nil sender")
			}
			receiver, err := user.ParseAddress(tx.GetReceiver())
			if err != nil {
				return nil, fmt.Errorf("nil receiver")
			}
			bal1, err := validator.bc.coll.Balance(sender)
			if err != nil {
				return nil, fmt.Errorf("something to do with balance sender")
			}
			bal2, err := validator.bc.coll.Balance(receiver)
			if err != nil {
				return nil, fmt.Errorf("something to do with balance receiver")
			}
			u1 := &user.User{
				Addr:    sender,
				Balance: bal1,
			}
			u2 := &user.User{
				Addr:    receiver,
				Balance: bal2,
			}
			return validator.bc.coll.AddUser(state, u1, u2), nil
		}
	}
	return nil, nil
}

func (validator *validator) validPostSize(post values.Bytes) bool {
	return len(post) > MinPostSize && len(post) < MaxPostSize
}

func (validator *validator) validMinReserveForPost(sender string) bool {
	addr, err := user.ParseAddress(sender)
	if err != nil {
		return false
	}
	bal, err := validator.bc.coll.Balance(addr)
	if err != nil {
		return false
	}
	return bal > uint64(MinTokenReserve)
}

func (validator *validator) validValue(sender string, value uint64) bool {
	addr, err := user.ParseAddress(sender)
	if err != nil {
		return false
	}
	bal, err := validator.bc.coll.Balance(addr)
	if err != nil {
		return false
	}
	return bal > value+value*uint64(Fees)/100
}
