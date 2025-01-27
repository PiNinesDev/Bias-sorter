package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", test)
	mux.HandleFunc("/test", index)

	log.Println("Server listening on port: 8080")
	http.ListenAndServe(":8080", mux)
	log.Println("done")
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
	log.Println("hit handler")
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").ParseFiles("index.html")
	if err != nil {
		log.Fatal("template error")
	}
	err = t.Execute(w, "test")
}
