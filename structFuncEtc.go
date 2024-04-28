package main

import (
	"database/sql"
	"encoding/json"
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

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{store: store}
}

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func wrappJsonError(message string) string {
	s, _ := json.Marshal(map[string]any{"error": message})
	return string(s)
}
