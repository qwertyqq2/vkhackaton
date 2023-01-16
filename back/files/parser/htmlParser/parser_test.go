package htmlparser

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	dataB, err := os.ReadFile("../../../htmlExample.html")
	if err != nil {
		log.Fatal(err)
	}
	data := string(dataB)
	comm, err := os.ReadFile("../../../comment.txt")
	if err != nil {
		log.Fatal(err)
	}
	newPost := NewParser(data).Add(string(comm))
	err = os.WriteFile("../../../htmlExample.html", []byte(newPost), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func TestParse(t *testing.T) {
	dataB, err := os.ReadFile("../../../htmlExample.html")
	if err != nil {
		log.Fatal(err)
	}
	data := string(dataB)
	body := NewParser(data).Body()
	head := NewParser(data).Head()
	fmt.Println("head:", head)
	fmt.Println("body: ", body)
}
