package main

import (
	"database/sql"
	"net/http"
)

type Server struct {
	DB     *sql.DB
	Router *Router
}

func newServer() *Server {
	s := &Server{}
	s.Router = &Router{}
	return s
}

func (s *Server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServerHTTP(w, r)
}
