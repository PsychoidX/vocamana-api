package test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetAllSentencesWithNoRows(t *testing.T) {
	// ログイン中のUserに紐づくSentenceが1つも無い場合nullが返ることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromSentences()

	// レコードが1つも無い場合、[]ではなくnullが返る
	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences",
		nil,
		nil,
		"",
		sc.GetAllSentences,
		http.StatusOK,
		"null",
	)
}

func TestGetAllSentences(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromSentences()

	var idWithUserId1 int
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'testsentence', 1)
		RETURNING id;
	`).Scan(&idWithUserId1)

	var idWithUserId2 int
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'testsentence', 2)
		RETURNING id;
	`).Scan(&idWithUserId2)

	expectedJSON := fmt.Sprintf(`
		[
			{
				"id": %d,
				"sentence": "testsentence",
				"user_id": 1
			}
		]`,
		idWithUserId1,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences",
		nil,
		nil,
		"",
		sc.GetAllSentences,
		http.StatusOK,
		expectedJSON,
	)
}

func TestGetSentenceByIdWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromSentences()

	var id string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'testsentence', 1)
		RETURNING id;
	`).Scan(&id)

	expectedJSON := fmt.Sprintf(`
		{
			"id": %s,
			"sentence": "testsentence",
			"user_id": 1
		}`,
		id,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences/:sentenceId",
		[]string{"sentenceId"},
		[]string{id},
		"",
		sc.GetSentenceById,
		http.StatusOK,
		expectedJSON,
	)
}

func TestGetSentenceByIdWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないSentenceを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromSentences()

	var id string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'testsentence', 2)
		RETURNING id;
	`).Scan(&id)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences/:sentenceId",
		[]string{"sentenceId"},
		[]string{id},
		"",
		sc.GetSentenceById,
		http.StatusOK,
		"{}",
	)
}

func TestCreateSentence(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを作成できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする

	// TODO
}

func TestUpdateSentence(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを更新できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ更新可能とする

	// TODO
}

func TestDeleteSentence(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを削除できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ削除可能とする

	// TODO
}

func TestAssociateSentenceWithWords(t *testing.T) {
	// WordとSentenceがどちらもログイン中のUserに紐づく場合
	// それらを紐づかせられることをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする

	// TODO
}