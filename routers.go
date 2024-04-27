package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./web")))

	r.Get("/api/nextdate", h.getNextDateHandler)
	return r
}
