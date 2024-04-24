package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/api/nextdate", getNextDateHandler)
	//http.HandleFunc("/api/nextdate", postTaskHandler)
	scheduler()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Println("Server is running on port localhost:7540")
	log.Fatal(http.ListenAndServe("localhost:7540", nil))
}
