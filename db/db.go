package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// ErrNoMatch is returned when we request a row that doesn't exist
var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.DB
}

func New(host, port, username, password, database string) (*Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return &db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return &db, err
	}
	log.Println("Database connection established")
	return &db, nil
}
