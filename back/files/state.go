package files

import "github.com/qwertyqq2/filebc/values"

type State interface {
	Add(values.Bytes, values.Bytes) values.Bytes

	Get(data ...values.Bytes) values.Bytes
}
