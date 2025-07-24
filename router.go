package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"example.com/bias-sorter/db"
)

type Router struct {
	Mux      http.ServeMux
	dataBase *sql.DB
	queries  *db.Queries
	sorter   Sorter
}

func newRouter(dataBase *sql.DB, staticDir string) *Router {
	r := &Router{}

	// returns the next or current pair of elements to compair
	r.Mux = *http.NewServeMux()
	r.queries = db.New(dataBase)
	r.dataBase = dataBase
	fs := http.FileServer(http.Dir(staticDir))

	r.Mux.Handle("/"+staticDir+"/", http.StripPrefix("/"+staticDir+"/", fs))
	r.Mux.HandleFunc("/quizes", makeViewQuizesHandler(r.queries))
	r.Mux.HandleFunc("/add-quiz", makeAddQuizHandler(r.queries))
	r.Mux.HandleFunc("/delete-quiz/{id}", makeDeleteQuizHandler(r.queries))
	r.Mux.HandleFunc("/add-entry/{id}", makeAddEntryHandler(r.queries))
	r.Mux.HandleFunc("/delete-entry/{id}", makeDeleteEntryHandler(r.queries))
	r.Mux.HandleFunc("/quiz/{quizID}", makeQuizHandler(r.queries))
	r.Mux.HandleFunc("/quiz/{quizID}/{sessionID}", makeTakeQuizHandler(r.queries))
	r.Mux.HandleFunc("/test", makeTestHandler(r.queries))

	log.Printf("router created\n")

	return r
}

func makeQuizHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("quiz handler called\n")
		quizID, err := strconv.Atoi(r.PathValue("quizID"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		ctx := context.Background()
		_, err = queries.GetQuiz(ctx, int64(quizID))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		sorter, err := queries.NewSorter(ctx, db.NewSorterParams{0, int64(quizID), sql.NullString{}})
		if err != nil {
			log.Fatalf("db error %v", err)
		}
		path := fmt.Sprintf("/quiz/%v/%v", sorter.QuizID, sorter.ID)
		http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	}
}

func makeTakeQuizHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("specific quiz handler called\n")
		sorterID, err := strconv.Atoi(r.PathValue("sessionID"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		ctx := context.Background()
		sorter, err := queries.GetSorter(ctx, int64(sorterID))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		t, err := template.ParseFiles("./templates/take-quiz.html")
		if err != nil {
			log.Fatalf("template error %v", err)
		}
		err = t.Execute(w, sorter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func makeDeleteEntryHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("delete entry handler called\n")
		entryID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatalf("id wasn't an int %v", err)
		}
		ctx := context.Background()
		queries.DeleteEntry(ctx, int64(entryID))
		log.Printf("delete entry %v\n", r.PathValue("id"))
		http.Redirect(w, r, "/quizes", http.StatusTemporaryRedirect)
	}
}

func makeAddEntryHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("add entry handler called\n")
		r.ParseForm()
		quizID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatalf("id wasn't an int %v", err)
		}
		entry := r.FormValue("entry-name")
		ctx := context.Background()
		queries.NewEntry(ctx, db.NewEntryParams{entry, int64(quizID)})
		http.Redirect(w, r, "/quizes", http.StatusTemporaryRedirect)
	}
}

func makeViewQuizesHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("view quiz handler called\n")
		ctx := context.Background()
		quizzes, err := queries.FindAllQuizes(ctx)
		var quizzesWithEntries []Quiz
		for _, quiz := range quizzes {
			entries, err := queries.GetQuizEntries(ctx, quiz.ID)
			if err != nil {
				log.Fatalf("query error %v", err)
			}
			quizzesWithEntries = append(quizzesWithEntries, Quiz{quiz.ID, quiz.Name, entries})
		}
		if err != nil {
			log.Fatalf("query error %v", err)
		}
		t, err := template.ParseFiles("./templates/list-quizzes.html")
		if err != nil {
			log.Fatalf("template error %v", err)
		}
		err = t.Execute(w, quizzesWithEntries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func makeAddQuizHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("add quiz handler called\n")
		r.ParseForm()
		ctx := context.Background()
		queries.NewQuiz(ctx, r.FormValue("quiz-name"))
		http.Redirect(w, r, "/quizes", http.StatusTemporaryRedirect)
	}
}

func makeDeleteQuizHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("delete quiz handler called\n")
		quizID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatalf("id wasn't an int %v", err)
		}
		ctx := context.Background()
		queries.DeleteQuiz(ctx, int64(quizID))
		log.Printf("delete quiz %v\n", r.PathValue("id"))
		http.Redirect(w, r, "/quizes", http.StatusTemporaryRedirect)
	}
}

func makeTestHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Test Handler Called\n")
		ctx := context.Background()
		quiz, _ := queries.NewQuiz(ctx, "Test_quiz")
		log.Printf("quiz added %v\n", quiz.Name)
		entry, _ := queries.NewEntry(ctx, db.NewEntryParams{"Test_entry", quiz.ID})
		log.Printf("entry added %v\n", entry.Name)
		quizEntries, _ := queries.GetQuizEntries(ctx, quiz.ID)
		log.Printf("quiz_entries %v\n", quizEntries)
		queries.DeleteQuiz(ctx, quiz.ID)
	}
}
