package main

import (
	"net/http"
)

type Router struct {
	Mux http.ServeMux
}

func (router *Router) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	return
}
