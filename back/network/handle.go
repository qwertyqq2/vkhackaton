package network

import (
	"net"
	"time"
)

func Handle(option int, conn Conn, pack *Package, handle func(*Package) string) error {
	if pack.Option != option {
		return ErrNotOpt
	}
	respPack := &Package{
		Option: option,
		Data:   handle(pack),
	}
	serializePack, err := respPack.SerializePack()
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(serializePack + EndBytes))
	return err
}

func Send(address string, pack *Package) *Package {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	defer conn.Close()
	serilizePack, err := pack.SerializePack()
	if err != nil {
		return nil
	}
	_, err = conn.Write([]byte(serilizePack + EndBytes))
	if err != nil {
		return nil
	}
	var (
		ch  = make(chan bool)
		res = new(Package)
	)
	go func() {
		res = readPack(conn)
		ch <- true
	}()
	select {
	case <-ch:

	case <-time.After(Waitime * time.Second):
		return nil
	}
	return res
}
