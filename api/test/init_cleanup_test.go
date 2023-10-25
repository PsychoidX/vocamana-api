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
	for i:=1; i<=2; i++ {
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

func DeleteAllFromWords() {
	// wordsテーブルのレコードを全件削除
	db.Exec("TRUNCATE TABLE words CASCADE;")
	// word_id_seqシーケンスを1にリセット
	// nextval()で、2から連番で取得される
	db.Exec("SELECT setval('word_id_seq', 1);")
}

func GetCurrentWordsSequenceValue() int {
	var currval int
	db.QueryRow(
		"SELECT currval('word_id_seq');",
	).Scan(&currval);
	return currval
}

func GetNextWordsSequenceValue() int {
	// インデックスのカウンタを進めず参照のみするための実装
	return GetCurrentWordsSequenceValue() + 1
}