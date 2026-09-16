package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("stub server listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", routes()))
}
