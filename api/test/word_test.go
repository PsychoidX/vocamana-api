package test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetAllWords(t *testing.T) {
	DeleteAllFromWords()

	// レコードが1つも無い場合、[]ではなくnullが返る
	DoSimpleTest(
		t,
		http.MethodGet,
		"/words",
		"",
		wc.GetAllWords,
		http.StatusOK,
		"null",
	)

	// レコードが存在する場合、ログイン中のユーザのレコードが全件返る
	// TODO

	// とりあえず、user_id=1のレコードだけ返すよう実装
	var idWithUserId1 int
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&idWithUserId1)

	var idWithUserId2 int
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword2', 'testmemo2', 2)
		RETURNING id;
	`).Scan(&idWithUserId2)

	expectedJSON := fmt.Sprintf(`
		[
			{
				"id": %d,
				"word": "testword",
				"memo": "testmemo",
				"user_id": 1
			}
		]`,
		idWithUserId1,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words",
		"",
		wc.GetAllWords,
		http.StatusOK,
		expectedJSON,
	)
}

func TestGetWordById(t *testing.T) {
	// TODO
}

func TestCreateWord(t *testing.T) {
	// TOOD
}

func TestUpdateWord(t *testing.T) {
	// TOOD
}

func TestDeleteWord(t *testing.T) {
	// TOOD
}
