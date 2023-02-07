package network

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

func handleStream(stream network.Stream) {
	log.Println("Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	closed := make(chan bool)

	go readData(rw, closed)
	go writeData(rw, closed, "boba")

}

func readData(rw *bufio.ReadWriter, closed chan bool) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Error reading from buffer")
			closed <- true
			break
		}
		if str == "" {
			return
		}
		if str != "\n" {
			log.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func writeData(rw *bufio.ReadWriter, closed chan bool, sendData string) {

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		_, err := rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			log.Println("Error writing to buffer")
			closed <- true
			break
		}
		err = rw.Flush()
		if err != nil {
			log.Println("Error flushing buffer")
			break
		}
	}
}
