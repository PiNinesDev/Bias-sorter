package main

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"example.com/bias-sorter/db"
	"golang.org/x/crypto/bcrypt"
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
	//r.Mux.HandleFunc("/quiz/{quizID}", makeQuizHandler(r.queries))
	//r.Mux.HandleFunc("/quiz/{quizID}/{sessionID}", makeTakeQuizHandler(r.queries))
	//r.Mux.HandleFunc("/test", makeTestHandler(r.queries))
	r.Mux.HandleFunc("POST /login", makePostLoginHandler(r.queries))
	r.Mux.HandleFunc("GET /login", makeGetLoginHandler(r.queries))
	r.Mux.HandleFunc("GET /signup", makeGetSignupHandler(r.queries))
	r.Mux.HandleFunc("POST /signup", makePostSignupHandler(r.queries))

	log.Printf("router created\n")

	return r
}

func makeGetSignupHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("signin GET handler called\n")

		t, err := template.ParseFiles("./templates/signup.html")
		if err != nil {
			log.Fatalf("template error %v", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func makePostSignupHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("signin POST handler called\n")
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Printf("username: %v, password: %v\n", username, password)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			log.Printf("issue with hashing")
			return
		}

		params := db.NewUserParams{
			Name:         username,
			PasswordHash: string(hashedPassword),
		}
		ctx := context.Background()
		user, err := queries.NewUser(ctx, params)
		if err != nil {
			http.Error(w, "Failed to create user in db", http.StatusInternalServerError)
			log.Printf("issue with db: %v", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		log.Printf("user create %v", user)

	}
}

func makeGetLoginHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("login GET handler called\n")
		t, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			log.Fatalf("template error %v", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func makePostLoginHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("login POST handler called\n")
		username := r.FormValue("username")
		password := r.FormValue("password")

		ctx := context.Background()
		user, err := queries.GetUserByUsername(ctx, username)
		if err != nil {
			http.Error(w, "Failed to find user in db", http.StatusInternalServerError)
			log.Printf("Failed to find user in db")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			http.Error(w, "Password did not match", http.StatusInternalServerError)
			log.Printf("Password did not match")
			return
		}

		w.WriteHeader(http.StatusAccepted)
		log.Printf("user verified %v", user)
	}
}

/*
func makeQuizHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("quiz handler called\n")
		quizID, err := strconv.Atoi(r.PathValue("quizID"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		ctx := context.Background()
		_, err = queries.GetQuizByID(ctx, int64(quizID))
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
*/

func makeDeleteEntryHandler(queries *db.Queries) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("delete entry handler called\n")
		entryID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatalf("id wasn't an int %v", err)
		}
		ctx := context.Background()
		queries.DeactivateEntry(ctx, int64(entryID))
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
		quizzes, err := queries.GetRectentQuizzes(ctx, db.GetRectentQuizzesParams{10, 1})
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
		params := db.NewQuizParams{
			Name:   r.FormValue("quiz-name"),
			UserID: 1, // TODO: change to get userID of logged in user
		}
		queries.NewQuiz(ctx, params)
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
		params := db.DeactivateQuizParams{
			ID:     int64(quizID),
			UserID: 1, // TODO: change to get userID of logged in user
		}
		queries.DeactivateQuiz(ctx, params)
		log.Printf("delete quiz %v\n", r.PathValue("id"))
		http.Redirect(w, r, "/quizes", http.StatusTemporaryRedirect)
	}
}
