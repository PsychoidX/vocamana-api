package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllWordsWithNoRows(t *testing.T) {
	// ログイン中のUserに紐づくWordが1つも無い場合nullが返ることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみUpdate可能とする
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
}

func TestGetAllWords(t *testing.T) {
	// ログイン中のUserに紐づくWordを取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみUpdate可能とする
	DeleteAllFromWords()

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

func TestGetWordByIdWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくWordを取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromWords()

	var id string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&id)

	// user_id=1の場合取得可能
	expectedJSON := fmt.Sprintf(`
		{
			"id": %s,
			"word": "testword",
			"memo": "testmemo",
			"user_id": 1
		}`,
		id,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{id},
		"",
		wc.GetWordById,
		http.StatusOK,
		expectedJSON,
	)
}

func TestGetWordByIdWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないWordを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromWords()

	var id string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword2', 'testmemo2', 2)
		RETURNING id;
	`).Scan(&id)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{id},
		"",
		wc.GetWordById,
		http.StatusOK,
		"{}",
	)
}

func TestCreateWord(t *testing.T) {
	// ログイン中のUserに紐づくWordを作成できることをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()

	nextId := GetNextWordsSequenceValue()

	reqBody := `{
		"word": "testword",
		"memo": "testmemo"
	}`

	// 登録されたレコードが返る
	expectedJSON := fmt.Sprintf(`
		{
			"id": %d,
			"word": "testword",
			"memo": "testmemo",
			"user_id": 1
		}`,
		nextId,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/words",
		nil,
		nil,
		reqBody,
		wc.CreateWord,
		http.StatusCreated,
		expectedJSON,
	)

	// DBにレコードが追加される
	var word string
	var memo string
	db.QueryRow(`
		SELECT word, memo FROM words
		WHERE id = $1;
	`,
	nextId,
	).Scan(&word, &memo)

	assert.Equal(t, "testword", word)
	assert.Equal(t, "testmemo", memo)
}

func TestUpdateWordWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくWordをUpdateできることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみUpdate可能とする
	DeleteAllFromWords()

	var id string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo', 1)
		RETURNING id;
	`).Scan(&id)

	reqBody := `{
		"word": "updated word",
		"memo": "updated memo"
	}`

	// 変更後のレコードが返る
	expectedJSON := fmt.Sprintf(`
		{
			"id": %s,
			"word": "updated word",
			"memo": "updated memo",
			"user_id": 1
		}`,
		id,
	)

	DoSimpleTest(
		t,
		http.MethodPut,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{id},
		reqBody,
		wc.UpdateWord,
		http.StatusAccepted,
		expectedJSON,
	)

	// DBのレコードが更新される
	var word string
	var memo string
	db.QueryRow(`
		SELECT word, memo FROM words
		WHERE id = $1;
	`,
	id,
	).Scan(&word, &memo)

	assert.Equal(t, "updated word", word)
	assert.Equal(t, "updated memo", memo)
}

func TestUpdateWordWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないWordをUpdateできないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみUpdate可能とする
	DeleteAllFromWords()

	var id string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo', 2)
		RETURNING id;
	`).Scan(&id)

	reqBody := `{
		"word": "updated word",
		"memo": "updated memo"
	}`

	DoSimpleTest(
		t,
		http.MethodPut,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{id},
		reqBody,
		wc.UpdateWord,
		http.StatusAccepted,
		"{}",
	)

	// DBのレコードが更新されない
	var word string
	var memo string
	db.QueryRow(`
		SELECT word, memo FROM words
		WHERE id = $1;
	`,
	id,
	).Scan(&word, &memo)

	assert.Equal(t, "word", word)
	assert.Equal(t, "memo", memo)
}

func TestDeleteWordWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくWordを削除できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみDelete可能とする
	DeleteAllFromWords()

	var id string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo', 1)
		RETURNING id;
	`).Scan(&id)

	// 削除したレコードが返る
	expectedJSON := fmt.Sprintf(`
		{
			"id": %s,
			"word": "word",
			"memo": "memo",
			"user_id": 1
		}`,
		id,
	)

	DoSimpleTest(
		t,
		http.MethodDelete,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{id},
		"",
		wc.DeleteWord,
		http.StatusAccepted,
		expectedJSON,
	)

	// DBからレコードが削除されている
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM words
		WHERE id = $1;
	`,
	id,
	).Scan(&count)

	assert.Equal(t, 0, count)
}

func TestDeleteWordWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないWordは削除できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみDelete可能とする
	DeleteAllFromWords()

	var id string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo', 2)
		RETURNING id;
	`).Scan(&id)

	DoSimpleTest(
		t,
		http.MethodDelete,
		"/words/:wordId",
		[]string{"wordId"},
		[]string{id},
		"",
		wc.DeleteWord,
		http.StatusAccepted,
		"{}",
	)

	// DBからレコードが削除されていない
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM words
		WHERE id = $1;
	`,
	id,
	).Scan(&count)

	assert.Equal(t, 1, count)
}