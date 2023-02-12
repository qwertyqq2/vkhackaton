package core

import (
	"encoding/binary"
	"fmt"

	"github.com/qwertyqq2/filebc/core/types"
	"github.com/qwertyqq2/filebc/core/types/transaction"
	"github.com/qwertyqq2/filebc/crypto"
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
	seed Seed

	v map[string]values.Bytes

	state values.Bytes
}

func newValidator(s Seed) *validator {
	return &validator{
		seed: s,
		v:    make(map[string]values.Bytes),
	}
}

func (validator *validator) add(state values.Bytes, txs ...types.Transaction) (values.Bytes, error) {

	var (
		sender      *user.Address
		receiver    *user.Address
		exSender    bool
		exReceiver  bool
		balSender   uint64
		balReceiver uint64
	)

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
			state = validator.seed.AddFile(state, file)

		case transaction.TypeTransferTx:

			if val, ok := validator.v[tx.GetSender()]; ok {
				balSender = binary.BigEndian.Uint64(val)
				send, err := user.ParseAddress(tx.GetSender())
				if err != nil {
					return nil, fmt.Errorf("pars addr err")
				}
				sender = send
				exSender = true
			} else {
				send, err := user.ParseAddress(tx.GetSender())
				if err != nil {
					return nil, fmt.Errorf("nil sender")
				}
				bal1, ex1, err := validator.seed.Balance(send)
				if err != nil {
					return nil, fmt.Errorf("something to do with balance sender")
				}
				exSender = ex1
				balSender = bal1
				sender = send
			}
			if val, ok := validator.v[tx.GetReceiver()]; ok {
				balReceiver = binary.BigEndian.Uint64(val)
				rec, err := user.ParseAddress(tx.GetReceiver())
				if err != nil {
					return nil, fmt.Errorf("pars addr err")
				}
				receiver = rec
				exReceiver = true
			} else {
				rec, err := user.ParseAddress(tx.GetReceiver())
				if err != nil {
					return nil, fmt.Errorf("nil sender")
				}
				bal2, ex2, err := validator.seed.Balance(rec)
				if err != nil {
					return nil, fmt.Errorf("something to do with balance sender")
				}
				exReceiver = ex2
				balReceiver = bal2
				receiver = rec
			}
			if !validator.validValue(balSender, tx.GetValue()) {
				return nil, fmt.Errorf("invalid value")
			}
			u1 := &user.User{
				Addr:    sender,
				Balance: balSender,
			}
			u2 := &user.User{
				Addr:    receiver,
				Balance: balReceiver,
			}

			if exSender && !exReceiver {
				invd := validator.seed.State().Inverse(u1.Hash())
				state = validator.seed.Add(state, invd)
				state = validator.seed.AddUser(state, user.GetUser(sender, balSender-tx.GetValue()),
					user.GetUser(receiver, balReceiver+tx.GetValue()))
			}
			if !exSender && exReceiver {
				invd := validator.seed.State().Inverse(u2.Hash())
				state = validator.seed.Add(state, invd)
				state = validator.seed.AddUser(state, user.GetUser(sender, balSender-tx.GetValue()),
					user.GetUser(receiver, balReceiver+tx.GetValue()))
			}
			if exSender && exReceiver {
				invn := validator.seed.State().Inverse(u1.Hash())
				invd := validator.seed.State().Inverse(u2.Hash())
				state = validator.seed.Add(state, invn)
				state = validator.seed.Add(state, invd)
				state = validator.seed.AddUser(state, user.GetUser(sender, balSender-tx.GetValue()),
					user.GetUser(receiver, balReceiver+tx.GetValue()))
			}
			if !exSender && !exReceiver {
				return nil, fmt.Errorf("users not exist")
			}

			validator.v[sender.String()] = crypto.ToBytes(balSender - tx.GetValue())
			validator.v[receiver.String()] = crypto.ToBytes(balSender + tx.GetValue())

		}
	}
	return state, nil
}

func (validator *validator) validPostSize(post values.Bytes) bool {
	return len(post) > MinPostSize && len(post) < MaxPostSize
}

func (validator *validator) validMinReserveForPost(sender string) bool {
	addr, err := user.ParseAddress(sender)
	if err != nil {
		return false
	}
	bal, _, err := validator.seed.Balance(addr)
	if err != nil {
		return false
	}
	return bal > uint64(MinTokenReserve)
}

func (validator *validator) validValue(bal uint64, value uint64) bool {
	return bal > value+value*uint64(Fees)/100
}
