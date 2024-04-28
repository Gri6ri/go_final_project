package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

const dateFormat = "20060102"

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

type Handler struct {
	store Store
}

func NewHandler(store Store) Handler {
	return Handler{store: store}
}

func wrappJsonError(message string) string {
	s, _ := json.Marshal(map[string]any{"error": message})
	return string(s)
}

func (h Handler) writeResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	_, err := w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
