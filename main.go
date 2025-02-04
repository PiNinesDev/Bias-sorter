package main

import (
	"html/template"
	"log"
	"net/http"
)

type Quiz struct {
	Title       string
	Description string
	Entries     []string
}

func main() {
	quizes := []Quiz{{"Best Super Smash Ultamite", "This quiz will rank all the fighters in Super Smash Ultamite", nil}}
	mux := http.NewServeMux()

	// set up static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// set up routes
	mux.HandleFunc("/home", makeIndexHandler(&quizes))
	mux.HandleFunc("/add-quiz", makeAddQuizHandler(&quizes))
	mux.HandleFunc("/add-entry", makeCreateQuizeAddEntry(&quizes))

	log.Println("Server listening on port: 8080")
	http.ListenAndServe(":8080", mux)
}

// TODO: replace q with database handler
func makeIndexHandler(q *[]Quiz) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			log.Fatal("template error")
		}

		err = t.Execute(w, q)
	}
}

// TODO: replace q with database handler and ad logic to choose the index via
// uri
func makeAddQuizHandler(q *[]Quiz) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle input
		t, err := template.ParseFiles("./templates/create-quiz.html", "./templates/partials/create-quiz-partials.html")
		if err != nil {
			log.Fatalf("template error: %v", err)
		}
		err = t.Execute(w, *q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("addQuiz template run")
	}
}

func makeCreateQuizeAddEntry(q *[]Quiz) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		newEntry := r.FormValue("new-entry")
		t, err := template.ParseFiles("./templates/partials/create-quiz-partials.html")

		if err != nil {
			log.Fatalf("template error: %v", err)
		}
		err = t.ExecuteTemplate(w, "createQuizEntry", newEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
