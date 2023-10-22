package test

import (
	"api/controller"
	"api/repository"
	"api/usecase"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // Postgresのドライバ
)

var db *sql.DB

// Word
var wr repository.IWordRepository
var wu *usecase.WordUsecase
var wc controller.IWordController

func TestMain(m *testing.M) {
	db = setupDB()

	wr = repository.NewWordRepository(db)
	wu = usecase.NewWordUsecase(wr)
	wc = controller.NewWordController(wu)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func setupDB() *sql.DB {
	pgConfigString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_POSTGRES_USER"),
		os.Getenv("TEST_POSTGRES_DB"),
		os.Getenv("TEST_POSTGRES_PASSWORD"),
	)
	db, err := sql.Open("postgres", pgConfigString)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func DeleteAllFromWords() {
	// wordsテーブルのレコードを削除し、シーケンスをリセットする
	db.Exec("TURNCATE TABLE words RESTART IDENTITY;")
}