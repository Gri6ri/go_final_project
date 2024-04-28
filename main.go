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
	r.Get("/api/nextdate", handler.getNextDateHandler)
	r.Get("/api/tasks", handler.getAllTasksHandler)
	r.Post("/api/task", handler.postTaskHandler)
	r.Get("/api/task", handler.getTaskHandler)
	r.Put("/api/task", handler.editTaskHandler)
	r.Delete("/api/task", handler.DeleteTaskHandler)
	r.Post("/api/task/done", handler.postTaskDoneHandler)

	log.Println("Server is running on port localhost:7540")
	log.Fatal(http.ListenAndServe("localhost:7540", nil))

}
