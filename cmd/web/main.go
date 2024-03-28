package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// var webPort = 80

func main() {
	// connect to database
	db := initDB()
	fmt.Println(db)
	db.Ping()

	// create sessions

	// create channels

	// create waitgroups

	// set up the application config

	// set up mails

	// listen for web connection
}

func initDB() *sql.DB {
	conn, err := connectToDB()

	if err != nil {
		log.Panic("cant connect to DB")
	}

	return conn
}

func connectToDB() (*sql.DB, error) {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Println("connected to database")
			return connection, nil
		}

		if counts > 10 {
			return nil, err
		}

		log.Print("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++

		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}