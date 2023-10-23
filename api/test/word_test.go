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
	db.Exec(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
	`)

	currval := GetCurrentWordsSequenceValue() // user_id=1のレコードが追加された時点でのid

	db.Exec(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword2', 'testmemo2', 2)
	`)

	expectedJSON := fmt.Sprintf(`
		[
			{
				"id": %d,
				"word": "testword",
				"memo": "testmemo",
				"user_id": 1
			}
		]`,
		currval,
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