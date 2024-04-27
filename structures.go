package main

import (
	"database/sql"
	"log"
	"net/http"
)

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

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{Addr: ":" + port, Handler: handler}
	log.Printf("Запуск сервера на порте: %s", port)
	return s.httpServer.ListenAndServe()
}
