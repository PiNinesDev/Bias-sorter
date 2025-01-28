package main

import (
	"html/template"
	"log"
	"net/http"
)

type Quiz struct {
	Title       string
	Description string
}

func main() {
	mux := http.NewServeMux()

	// set up static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// set up routes
	mux.HandleFunc("/", index)
	mux.HandleFunc("/add-quiz", index)

	log.Println("Server listening on port: 8080")
	http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal("template error")
	}

	quizes := []Quiz{{"Best Super Smash Ultamite", "This quiz will rank all the fighters in Super Smash Ultamite"}}
	err = t.Execute(w, quizes)
}

func addQuiz(w http.ResponseWriter, r *http.Request) {
	// TODO handle post request

}
