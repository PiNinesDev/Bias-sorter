package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	log.Println("Server listening on port: 8080")
	http.ListenAndServe(":8080", mux)
	log.Println("done")
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal("template error")
	}

	type testArt struct {
		Title string
		Desc  string
	}

	articals := []testArt{{"title", "This is a longer description"}}
	err = t.Execute(w, articals)
}
