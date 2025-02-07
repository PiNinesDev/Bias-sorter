package main

import (
	"crypto"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Quiz struct {
	Title       string
	Description string
	Entries     []string
}

func main() {
	db := initDB()
	defer db.Close()

	quizes := []Quiz{{"Best Super Smash Ultamite", "This quiz will rank all the fighters in Super Smash Ultamite", nil}}
	mux := http.NewServeMux()

	// set up static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// set up routes
	mux.HandleFunc("/home", makeIndexHandler(&quizes))
	mux.HandleFunc("/create-quiz", makeAddQuizHandler(&quizes))
	mux.HandleFunc("/add-entry", makeCreateQuizeAddEntry(&quizes))
	mux.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/add-quiz", makeCreateQuiz(db))

	log.Println("Server listening on port: 8080")
	http.ListenAndServe(":8080", mux)
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	createQuizTable := `
	CREATE TABLE IF NOT EXISTS quizzes (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		created_on DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createQuizTable)
	if err != nil {
		log.Println(err)
	}

	return db
}

func addQuizToDB(db *sql.DB, title string, description string) error {
	query := `
	INSERT INTO quizzes (id, title, description) VALUES (?,?,?)
	`
	id := uuid.New().String()
	_, err := db.Exec(query, id, title, description)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func makeCreateQuiz(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println(r.Form)
		addQuizToDB(db, r.FormValue("QuizName"), r.FormValue("QuizDesc"))
	}
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
