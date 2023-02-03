package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	data1, err := os.ReadFile("htmlfiles/htmlExample1.html")
	if err != nil {
		log.Fatal(err)
	}
	data2, err := os.ReadFile("htmlfiles/htmlExample2.html")
	if err != nil {
		log.Fatal(err)
	}
	data3, err := os.ReadFile("htmlfiles/htmlExample3.html")
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.Join(
		[][]byte{
			data1, data2, data3,
		},
		[]byte{},
	)
	_, err = w.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func TestHandle(t *testing.T) {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}
