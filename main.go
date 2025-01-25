package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", test)

	log.Println("Server listening on port: 8080")
	http.ListenAndServe(":8080", mux)
	log.Println("done")
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
	log.Println("hit handler")
}
