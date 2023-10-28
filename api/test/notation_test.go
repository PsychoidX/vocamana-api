package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNotationWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを作成できることをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&wordId)

	id := GetNextNotationsSequenceValue()

	reqBody := `{
		"notation": "testnotation"
	}`

	// 登録されたレコードが返る
	resJSON := fmt.Sprintf(`
		{
			"id": %d,
			"word_id": %s,
			"notation": "testnotation"
		}`,
		id,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordId},
		reqBody,
		nc.CreateNotation,
		http.StatusCreated,
		resJSON,
	)

	// DBにレコードが追加される
	var notation string
	db.QueryRow(`
		SELECT notation FROM notations
		WHERE id = $1
			AND word_id = $2;
	`,
	id,
	wordId,
	).Scan(&notation)

	assert.Equal(t, "testnotation", notation)
}

func TestCreateNotationWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないWordに対し、Notationを作成できないことをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 2)
		RETURNING id;
	`).Scan(&wordId)

	id := GetNextNotationsSequenceValue()

	reqBody := `{
		"notation": "testnotation"
	}`

	DoSimpleTest(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordId},
		reqBody,
		nc.CreateNotation,
		http.StatusUnauthorized,
		"{}",
	)

	// DBにレコードが追加されない
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM notations
		WHERE id = $1
			AND word_id = $2;
	`,
	id,
	wordId,
	).Scan(&count)

	assert.Equal(t, 0, count)
}