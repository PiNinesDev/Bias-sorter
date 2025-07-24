package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"log"
)

func main() {
	database := initDB()
	defer database.Close()

	s := newServer(database)
	s.Start()
}

func initDB() *sql.DB {
	dbPath := "test.db"

	// open db
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// init goose
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}
	migrationsPath := "./migrations"
	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatal(err)

	}

	log.Println("migrations run succeful")
	return db
}
