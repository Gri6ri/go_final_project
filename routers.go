package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (s Service) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./web")))

	r.Get("/api/nextdate", s.getNextDateHandler)
	return r
}
