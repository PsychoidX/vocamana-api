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

// SentencesWords
var swr repository.ISentencesWordsRepository

// Sentence
var sr repository.ISentenceRepository
var su *usecase.SentenceUsecase
var sc controller.ISentenceController

// Notation
var nr repository.INotationRepository
var nu *usecase.NotationUsecase
var nc controller.INotationController

func TestMain(m *testing.M) {
	db = setupDB()

	// Repository
	wr = repository.NewWordRepository(db)
	sr = repository.NewSentenceRepository(db)
	swr = repository.NewSentencesWordsRepository(db)
	nr = repository.NewNotationRepository(db)

	// Usecase
	wu = usecase.NewWordUsecase(wr, sr, swr, nr)
	su = usecase.NewSentenceUsecase(sr, wr, swr, nr)
	nu = usecase.NewNotationUsecase(wr, sr, swr, nr)

	// Controller
	wc = controller.NewWordController(wu)
	sc = controller.NewSentenceController(su)
	nc = controller.NewNotationController(nu)

	setupUserData()

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

func setupUserData() {
	// テストのため、id=1, 2のユーザが存在しない場合に作成する
	for i := 1; i <= 2; i++ {
		db.Exec(`
			INSERT INTO users
			(id, email, password)
			SELECT CAST($1 AS INTEGER), CONCAT('sample', $1, '@example.com'), 'pass'
			WHERE NOT EXISTS(
				SELECT 1
				FROM users
				WHERE id = CAST($1 AS INTEGER)
			);`, i)
	}
}