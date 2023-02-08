package network

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
	for {
		data, err := ioutil.ReadAll(rw)
		if err != nil {
			log.Println("Error reading from buffer")
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
