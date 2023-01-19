package files

type State interface {
	Add([]byte, []byte) []byte

	Get(data ...[]byte) []byte
}
