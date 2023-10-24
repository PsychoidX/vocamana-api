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
		nil,
		nil,
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
		nil,
		nil,
		"",
		wc.GetAllWords,
		http.StatusOK,
		expectedJSON,
	)
}

func TestGetWordById(t *testing.T) {
	DeleteAllFromWords()

	// ログイン中のUserに紐づくWordだけを取得可能
	// TODO
	// とりあえず、user_id=1のUserに紐づくWordだけを取得可能なよう実装

	var idWithUserId1 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&idWithUserId1)

	// user_id=1の場合取得可能
	expectedJSON := fmt.Sprintf(`
		{
			"id": %s,
			"word": "testword",
			"memo": "testmemo",
			"user_id": 1
		}`,
		idWithUserId1,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{idWithUserId1},
		"",
		wc.GetWordById,
		http.StatusOK,
		expectedJSON,
	)

	// user_id=2の場合{}が返る
	var idWithUserId2 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword2', 'testmemo2', 2)
		RETURNING id;
	`).Scan(&idWithUserId2)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{idWithUserId2},
		"",
		wc.GetWordById,
		http.StatusOK,
		"{}",
	)
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
