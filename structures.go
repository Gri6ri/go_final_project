package main

import "database/sql"

type Task struct {
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

func newHandler(service Service) Handler {
	return Handler{service: service}
}
