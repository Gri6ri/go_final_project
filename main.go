package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbFile := getDbfile()

	db := getDb(dbFile)

	store := NewStore(db)

	service := NewService(store)
	http.HandleFunc("/api/nextdate", service.getNextDateHandler)
	// r := chi.NewRouter()

	// r.Get("/api/nextdate", service.getNextDateHandler)

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Println("Server is running on port localhost:7540")
	log.Fatal(http.ListenAndServe("localhost:7540", nil))

}
