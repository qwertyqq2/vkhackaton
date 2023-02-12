package files

import "github.com/qwertyqq2/filebc/values"

type State interface {
	Add(values.Bytes, ...values.Bytes) values.Bytes

	Get(...values.Bytes) values.Bytes

	Inverse(values.Bytes) values.Bytes
}
