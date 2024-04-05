package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Print("Server is running on the port localhost:7540")
	log.Fatal(http.ListenAndServe("localhost:7540", nil))
}
