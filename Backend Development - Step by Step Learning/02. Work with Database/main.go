package main

import (
	"LearnGoDB/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Need:
/*
	1. data source name (here 'dsn')
	2. Database Library (To connect database)
	   We are using Postgres.
	   So,
	   we need postgress implementation of that single api
	   that go provides a single interface for RDBMS and that
	   would enable us communicate with postgres
	   So,
	   we need to install that library: go get github.com/lib/pq
*/
func connectToDB(dsn string) (*sql.DB, error) {

	// Open connection to database
	db, err := sql.Open("postgres", dsn) // Open(driver name, data source name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type application struct {
	Models models.Models
}

func main() {

	// Connect with Postgres Database
	dsn := "user=postgres dbname=GoDB password=101010 sslmode=disable"

	db, err := connectToDB(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	app := application{
		Models: models.NewModel(db),
	}

	fmt.Println("Starting application...")
	err = app.serve()
	if err != nil {
		log.Fatalln(err)
	}
}