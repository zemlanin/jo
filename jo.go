package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	go h.run()
	http.HandleFunc("/", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}