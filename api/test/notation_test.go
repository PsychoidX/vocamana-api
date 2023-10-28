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
	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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

func TestUpdateNotation(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを更新できることをテスト
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

	var notationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation')
		RETURNING id;
	`,
	wordId,
	).Scan(&notationId)

	reqBody := `{
		"notation": "updated notation"
	}`

	// 更新されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %s,
			"word_id": %s,
			"notation": "updated notation"
		}`,
		notationId,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodPut,
		"/words/:wordId/notations/:notationId",
		[]string{"wordId", "notationId"},
		[]string{wordId, notationId},
		reqBody,
		nc.UpdateNotation,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBのレコードが更新される
	var notation string
	db.QueryRow(`
		SELECT notation FROM notations
		WHERE id = $1
			AND word_id = $2;
	`,
	notationId,
	wordId,
	).Scan(&notation)

	assert.Equal(t, "updated notation", notation)
}

func TestUpdateNotationWithNoRows(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、更新対象のNotationが無い場合、{}が返ることをテスト
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

	reqBody := `{
		"notation": "updated notation"
	}`

	DoSimpleTest(
		t,
		http.MethodPut,
		"/words/:wordId/notations/:notationId",
		[]string{"wordId", "notationId"},
		[]string{wordId, "1"},
		reqBody,
		nc.UpdateNotation,
		http.StatusUnauthorized,
		"{}",
	)
}

func TestDeleteNotationWithLoggingIn(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを削除できることをテスト
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

	var notationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation')
		RETURNING id;
	`,
	wordId,
	).Scan(&notationId)

	// 削除されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %s,
			"word_id": %s,
			"notation": "test notation"
		}`,
		notationId,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodDelete,
		"/words/:wordId/notations/:notationId",
		[]string{"wordId", "notationId"},
		[]string{wordId, notationId},
		"",
		nc.DeleteNotation,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBのレコードが削除される
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM notations
		WHERE id = $1
			AND word_id = $2;
	`,
	notationId,
	wordId,
	).Scan(&count)

	assert.Equal(t, 0, count)
}

func TestDeleteNotationWithoutLoggingIn(t *testing.T) {
	// ログイン中のUserに紐づかないWordに対し、Notationを削除できないことをテスト
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

	var notationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('word_id_seq'), $1, 'test notation')
		RETURNING id;
	`,
	wordId,
	).Scan(&notationId)

	DoSimpleTest(
		t,
		http.MethodDelete,
		"/words/:wordId/notations/:notationId",
		[]string{"wordId", "notationId"},
		[]string{wordId, notationId},
		"",
		nc.DeleteNotation,
		http.StatusUnauthorized,
		"{}",
	)

	// DBのレコードが削除される
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM notations
		WHERE id = $1
			AND word_id = $2;
	`,
	notationId,
	wordId,
	).Scan(&count)

	assert.Equal(t, 1, count)
}