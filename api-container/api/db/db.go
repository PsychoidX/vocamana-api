package db

import (
	"database/sql"
  _	"github.com/lib/pq" // Postgresのドライバ
	"fmt"
	"log"
	"os"
)

func NewDB() *sql.DB {
	pgConfigString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	)
	db, err := sql.Open("postgres", pgConfigString)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}