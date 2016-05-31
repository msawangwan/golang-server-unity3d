package main

import (
	"io"
	"log"
	"net/http"
)

var SERVER_LISTENING_PORT = ":8080"

func main() {
	log.Println("Server Started ... ")

	http.HandleFunc("/", handleTestFromUnity)
	http.ListenAndServe(SERVER_LISTENING_PORT, nil)
}

func handleTestFromUnity(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!")
}
