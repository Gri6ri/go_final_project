package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbFile := getDbfile()
	CreationOfDb := isDbHere(dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if CreationOfDb {
		createDb(dbFile, db)
	}
	store := NewStore(db)
	service := NewService(store)
	handler := NewHandler(service)

	http.Handle("/", http.FileServer(http.Dir("./web")))
	r := chi.NewRouter()
	http.Handle("/api/", r)
	r.HandleFunc("/api/nextdate", handler.getNextDateHandler)
	r.HandleFunc("/api/task", handler.postTaskHandler)

	log.Println("Server is running on port localhost:7540")
	log.Fatal(http.ListenAndServe("localhost:7540", nil))

}
