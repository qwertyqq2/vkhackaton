package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("htmlExample.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(data))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}
