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
	// とりあえずuser_id=1のWordのみ取得可能とする
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
	// とりあえずuser_id=1のWordのみ取得可能とする
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

	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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

	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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
	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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
	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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
		http.StatusUnauthorized,
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
	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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
		http.StatusUnauthorized,
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

func TestGetAssociatedSentences(t *testing.T) {
	// WordとSentenceがどちらもログイン中のuser_idに紐づく場合、
	// Wordに紐づくSentenceを取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word', 'test memo', 1)
		RETURNING id;
	`).Scan(&wordId)

	var sentenceId1 string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence1', 1)
		RETURNING id;
	`).Scan(&sentenceId1)

	var sentenceId2 string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence2', 1)
		RETURNING id;
	`).Scan(&sentenceId2)

	db.QueryRow(`
		INSERT INTO sentences_words
		(word_id, sentence_id)
		VALUES
		($1, $2),
		($1, $3);
		`,
		wordId,
		sentenceId1,
		sentenceId2,
	)

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %s,
				"sentence": "test sentence1",
				"user_id": 1
			},
			{
				"id": %s,
				"sentence": "test sentence2",
				"user_id": 1
			}
		]`,
		sentenceId1,
		sentenceId2,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/associated-sentences",
		[]string{"wordId"},
		[]string{wordId},
		"",
		wc.GetAssociatedSentences,
		http.StatusOK,
		expectedResponse,
	)
}

func TestGetAssociatedSentencesWithInvalidWordId(t *testing.T) {
	// Wordがログイン中のuser_idに紐づかない場合、
	// Wordに紐づくSentenceを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	// Wordのuser_idとしてログイン中のuser_id以外を使用
	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word', 'test memo', 2)
		RETURNING id;
	`).Scan(&wordId)

	var sentenceId1 string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence1', 1)
		RETURNING id;
	`).Scan(&sentenceId1)

	var sentenceId2 string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence2', 1)
		RETURNING id;
	`).Scan(&sentenceId2)

	db.QueryRow(`
		INSERT INTO sentences_words
		(word_id, sentence_id)
		VALUES
		($1, $2),
		($1, $3);
		`,
		wordId,
		sentenceId1,
		sentenceId2,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/associated-sentences",
		[]string{"wordId"},
		[]string{wordId},
		"",
		wc.GetAssociatedSentences,
		http.StatusOK,
		"null",
	)
}

func TestGetAssociatedSentencesWithInvalidSentenceId(t *testing.T) {
	// Sentenceがログイン中のuser_idに紐づかない場合、
	// Wordに紐づくSentenceを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word', 'test memo', 1)
		RETURNING id;
	`).Scan(&wordId)

	// Sentenceのuser_idとしてログイン中のuser_id以外を使用
	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence', 2)
		RETURNING id;
	`).Scan(&sentenceId)

	db.QueryRow(`
		INSERT INTO sentences_words
		(word_id, sentence_id)
		VALUES
		($1, $2);
		`,
		wordId,
		sentenceId,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/associated-sentences",
		[]string{"wordId"},
		[]string{wordId},
		"",
		wc.GetAssociatedSentences,
		http.StatusOK,
		"null",
	)
}

func TestGetAssociatedSentencesWithLink(t *testing.T) {
	// WordとSentenceがどちらもログイン中のuser_idに紐づく場合、
	// Wordに紐づくSentenceをリンク付きで取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'text word notation 1 text word notation 2 text', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'word memo', 1)
		RETURNING id;
	`).Scan(&wordId)

	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES
		(nextval('word_id_seq'), $1, 'word notation 1'),
		(nextval('word_id_seq'), $1, 'word notation 2');
		`,
		wordId,
	)

	db.QueryRow(`
		INSERT INTO sentences_words
		(sentence_id, word_id)
		VALUES
		($1, $2);
		`,
		sentenceId,
		wordId,
	)

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %s,
				"sentence": "text <a href=\"words/%s\">word notation 1</a> text <a href=\"words/%s\">word notation 2</a> text",
				"user_id": 1
			}
		]`,
		sentenceId,
		wordId,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/associated-sentences",
		[]string{"wordId"},
		[]string{wordId},
		"",
		wc.GetAssociatedSentencesWithLink,
		http.StatusOK,
		expectedResponse,
	)
}
