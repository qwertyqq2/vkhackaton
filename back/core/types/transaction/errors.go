package transaction

import "errors"

const (
	TypePostTx     = 101
	TypeTransferTx = 102
)

var (
	ErrIncorrectPost = errors.New("incorrect post")
)
