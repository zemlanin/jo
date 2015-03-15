package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	go h.run()
	http.HandleFunc("/", serveWs)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
