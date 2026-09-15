package main

import (
	"log"
	"net/http"
)

func main() {
	s := newServer()
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", s.routes()))
}
