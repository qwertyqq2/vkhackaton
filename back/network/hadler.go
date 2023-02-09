package network

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
)

func handleStream(stream network.Stream) {
	log.Println("Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	closed := make(chan bool)

	go readData(rw, closed)
	go writeData(rw, closed)

}

func readData(rw *bufio.ReadWriter, closed chan bool) {
	buf := make([]byte, 40*1024)
	for {
		n, err := rw.Read(buf)
		if err != nil {
			log.Println("Error reading from buffer")
			closed <- true
			break
		}
		data := buf[:n]
		if data == nil {
			log.Println("nil data")
			closed <- true
			break
		}
		msg, err := Unmarhsal(data)
		if err != nil {
			log.Println(err)
			closed <- true
			return
		}
		if msg == nil {
			log.Println("nil msg")
			closed <- true
			return
		}
		fmt.Println("out: ", string(msg.payload))
	}
}

func writeData(rw *bufio.ReadWriter, closed chan bool) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from stdin")
			break
		}
		msg := NewMessage(MsgName1, []byte(sendData))
		data, err := Marhal(msg)
		if err != nil {
			log.Println(err)
			break
		}
		if data == nil {
			log.Println("nil data")
			break
		}
		_, err = rw.Write(data)
		if err != nil {
			log.Println("Error writing to buffer")
			break
		}
		err = rw.Flush()
		if err != nil {
			log.Println("Error flushing buffer")
			break
		}
	}
}
