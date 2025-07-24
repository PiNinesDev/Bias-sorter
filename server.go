package main

import (
	"database/sql"
	"net/http"

	"example.com/bias-sorter/db"
)

type Server struct {
	DB     *sql.DB
	Router *Router
}

type Quiz struct {
	ID      int64
	Name    string
	Entries []db.Entry
}

func newServer(db *sql.DB) *Server {
	s := &Server{}

	s.Router = newRouter(db, "static")
	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":8080", &s.Router.Mux)
}
