package db

import (
	"database/sql"
  _	"github.com/lib/pq" // Postgresのドライバ
	"fmt"
	"log"
	"os"
)

func NewDB() *sql.DB {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname:= os.Getenv("POSTGRES_DB")
	args := fmt.Sprintf("user=%s dbname=%s password=%s", username, dbname, password)
	db, err := sql.Open("postgres", args)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}