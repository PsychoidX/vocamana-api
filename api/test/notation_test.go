package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllNotations(t *testing.T) {
	// Wordがログイン中のUserに紐づく場合、Wordに紐づくNotationを取得できることをテスト
	// TODO ログイン機能
	// とりあえずWordに紐づくUserがuser_id=1の場合のみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&wordId)

	var notationId1 string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation1')
		RETURNING id;
	`,
	wordId,
	).Scan(&notationId1)

	var notationId2 string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation2')
		RETURNING id;
	`,
	wordId,
	).Scan(&notationId2)

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %s,
				"word_id": %s,
				"notation": "test notation1"
			},
			{
				"id": %s,
				"word_id": %s,
				"notation": "test notation2"
			}
		]`,
		notationId1,
		wordId,
		notationId2,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordId},
		"",
		nc.GetAllNotations,
		http.StatusOK,
		expectedResponse,
	)
}

func TestGetAllNotationsWithNoRows(t *testing.T) {
	// Wordがログイン中のUserに紐づき、
	// かつWordに紐づくNotationの数が0の場合、nullが返ることをテスト
	// TODO ログイン機能
	// とりあえずWordに紐づくUserがuser_id=1の場合のみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&wordId)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordId},
		"",
		nc.GetAllNotations,
		http.StatusOK,
		"null",
	)
}

func TestGetAllNotationsWithInvalidWordId(t *testing.T) {
	// Wordがログイン中のUserに紐づかない場合、Wordに紐づくNotationを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずWordに紐づくUserがuser_id=1の場合のみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	var wordIdWithUserId1 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 1)
		RETURNING id;
	`).Scan(&wordIdWithUserId1)

	var wordIdWithUserId2 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'testword', 'testmemo', 2)
		RETURNING id;
	`).Scan(&wordIdWithUserId2)

	var notationIdWithUserId1 string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation1')
		RETURNING id;
	`,
	wordIdWithUserId1,
	).Scan(&notationIdWithUserId1)

	var notationIdWithUserId2 string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation2')
		RETURNING id;
	`,
	wordIdWithUserId2,
	).Scan(&notationIdWithUserId2)

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %s,
				"word_id": %s,
				"notation": "test notation1"
			}
		]`,
		notationIdWithUserId1,
		wordIdWithUserId1,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordIdWithUserId1},
		"",
		nc.GetAllNotations,
		http.StatusOK,
		expectedResponse,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordIdWithUserId2},
		"",
		nc.GetAllNotations,
		http.StatusOK,
		"null",
	)
}

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