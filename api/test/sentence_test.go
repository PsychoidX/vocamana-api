package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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

	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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

	expectedResponse := fmt.Sprintf(`
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
		expectedResponse,
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
	DeleteAllFromSentences()

	nextId := GetNextSentencesSequenceValue()

	reqBody := `{
		"sentence": "testsentence"
	}`

	// 登録されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %d,
			"sentence": "testsentence",
			"user_id": 1
		}`,
		nextId,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences",
		nil,
		nil,
		reqBody,
		sc.CreateSentence,
		http.StatusCreated,
		expectedResponse,
	)

	// DBにレコードが追加される
	var sentence string
	db.QueryRow(`
		SELECT sentence FROM sentences
		WHERE id = $1;
	`,
		nextId,
	).Scan(&sentence)

	assert.Equal(t, "testsentence", sentence)
}

func TestCreateSentenceIncludingWords(t *testing.T) {
	// 登録済みのWordの中に、新規追加されたSentence中に含まれるものがある場合、
	// sentences_wordsに追加されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var appleWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'りんご', 'test memo', 1)
		RETURNING id;
	`).Scan(&appleWordId)

	var redWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), '赤い', 'test memo', 1)
		RETURNING id;
	`).Scan(&redWordId)

	var blueWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), '青い', 'test memo', 1)
		RETURNING id;
	`).Scan(&blueWordId)

	sentenceId := GetNextSentencesSequenceValue()

	reqBody := `{
		"sentence": "赤いりんごを食べた"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/sentences",
		nil,
		nil,
		reqBody,
		sc.CreateSentence,
	)

	// 「赤いりんごを食べた」には「赤い」が含まれるため、
	// sentences_wordsに追加される
	var redCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		redWordId,
	).Scan(&redCount)

	assert.Equal(t, 1, redCount)

	// 「赤いりんごを食べた」には「りんご」が含まれるため、
	// sentences_wordsに追加される
	var appleCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		appleWordId,
	).Scan(&appleCount)

	assert.Equal(t, 1, appleCount)

	// 「赤いりんごを食べた」には「青い」が含まれないため、
	// sentences_wordsに追加されない
	var blueCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		blueWordId,
	).Scan(&blueCount)

	assert.Equal(t, 0, blueCount)
}

func TestCreateSentenceIncludingInvalidWords(t *testing.T) {
	// 登録済みのWordの中に、新規追加されたSentence中に含まれるものがあるが、
	// 含まれるWordのUserIdがログイン中のuser_idと異なる場合
	// sentences_wordsに追加されないことをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'りんご', 'test memo', 2)
		RETURNING id;
	`).Scan(&wordId)

	sentenceId := GetNextSentencesSequenceValue()

	reqBody := `{
		"sentence": "赤いりんごを食べた"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/sentences",
		nil,
		nil,
		reqBody,
		sc.CreateSentence,
	)

	// 「赤いりんごを食べた」には「りんご」が含まれるが、
	// Wordのuser_idがログイン中のものと異なるため、
	// sentences_wordsに追加されない
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		wordId,
	).Scan(&count)

	assert.Equal(t, 0, count)
}

func TestCreateSentenceIncludingNotations(t *testing.T) {
	// 登録済みのNotationの中に、新規追加されたSentence中に含まれるものがある場合、
	// sentences_wordsに追加されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var appleWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'りんご', 'test memo', 1)
		RETURNING id;
	`).Scan(&appleWordId)

	// 「りんご」の別表記として「林檎」を追加
	var appleNotationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('notation_id_seq'), $1, '林檎')
		RETURNING id;
	`,
		appleWordId,
	).Scan(&appleNotationId)

	var eatWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), '食べる', 'test memo', 1)
		RETURNING id;
	`).Scan(&eatWordId)

	// 「食べる」の別表記として「食う」を追加
	var eatNotationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('notation_id_seq'), $1, '食う')
		RETURNING id;
	`,
		eatWordId,
	).Scan(&eatNotationId)

	var fruitWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), '果物', 'test memo', 1)
		RETURNING id;
	`).Scan(&fruitWordId)

	// 「果物」の別表記として「果実」を追加
	var fruitNotationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('notation_id_seq'), $1, '果実')
		RETURNING id;
	`,
		fruitWordId,
	).Scan(&fruitNotationId)

	sentenceId := GetNextSentencesSequenceValue()

	reqBody := `{
		"sentence": "赤い林檎を食う"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/sentences",
		nil,
		nil,
		reqBody,
		sc.CreateSentence,
	)

	// 「赤い林檎を食う」には「りんご」の別表記「林檎」が含まれるため、
	// sentences_wordsに追加される
	var appleCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		appleWordId,
	).Scan(&appleCount)

	assert.Equal(t, 1, appleCount)

	// 「赤い林檎を食う」には「食べる」の別表記「食う」が含まれるため、
	// sentences_wordsに追加される
	var eatCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		eatWordId,
	).Scan(&eatCount)

	assert.Equal(t, 1, eatCount)

	// 「赤い林檎を食う」には「果物」も、その別表記「果実」も含まれないため、
	// sentences_wordsに追加されない
	var fruitCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		fruitWordId,
	).Scan(&fruitCount)

	assert.Equal(t, 0, fruitCount)
}

func TestCreateSentenceIncludingInvalidNotations(t *testing.T) {
	// 登録済みのNotationの中に、新規追加されたSentence中に含まれるものがあるが、
	// 含まれるNotationが紐づくWordのUserIdがログイン中のuser_idと異なる場合
	// sentences_wordsに追加されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var appleWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'りんご', 'test memo', 2)
		RETURNING id;
	`).Scan(&appleWordId)

	// 「りんご」の別表記として「林檎」を追加
	var appleNotationId string
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('notation_id_seq'), $1, '林檎')
		RETURNING id;
	`,
		appleWordId,
	).Scan(&appleNotationId)

	sentenceId := GetNextSentencesSequenceValue()

	reqBody := `{
		"sentence": "赤い林檎を食う"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/sentences",
		nil,
		nil,
		reqBody,
		sc.CreateSentence,
	)

	// 「赤い林檎を食う」には「りんご」の別表記「林檎」が含まれるが、
	// Wordのuser_idがログイン中のものと異なるため、
	// sentences_wordsに追加されない
	var appleCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		appleWordId,
	).Scan(&appleCount)

	assert.Equal(t, 0, appleCount)
}

func TestUpdateSentenceWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを更新できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ更新可能とする
	DeleteAllFromSentences()

	var id string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&id)

	reqBody := `{
		"sentence": "updated sentence"
	}`

	// 変更後のレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %s,
			"sentence": "updated sentence",
			"user_id": 1
		}`,
		id,
	)

	DoSimpleTest(
		t,
		http.MethodPut,
		"/words/:sentenceId",
		[]string{"sentenceId"},
		[]string{id},
		reqBody,
		sc.UpdateSentence,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBのレコードが更新される
	var sentence string
	db.QueryRow(`
		SELECT sentence FROM sentences
		WHERE id = $1;
	`,
		id,
	).Scan(&sentence)

	assert.Equal(t, "updated sentence", sentence)
}

func TestUpdateSentenceWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないSentenceを更新できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ更新可能とする
	DeleteAllFromSentences()

	var id string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 2)
		RETURNING id;
	`).Scan(&id)

	reqBody := `{
		"sentence": "updated sentence"
	}`

	DoSimpleTest(
		t,
		http.MethodPut,
		"/words/:sentenceId",
		[]string{"sentenceId"},
		[]string{id},
		reqBody,
		sc.UpdateSentence,
		http.StatusUnauthorized,
		"{}",
	)

	// DBのレコードが更新されない
	var sentence string
	db.QueryRow(`
		SELECT sentence FROM sentences
		WHERE id = $1;
	`,
		id,
	).Scan(&sentence)

	assert.Equal(t, "sentence", sentence)
}

func TestDeleteSentenceWithLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを削除できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ削除可能とする
	DeleteAllFromSentences()

	var id string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&id)

	// 削除したレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %s,
			"sentence": "sentence",
			"user_id": 1
		}`,
		id,
	)

	DoSimpleTest(
		t,
		http.MethodDelete,
		"/sentences/:sentenceId",
		[]string{"sentenceId"},
		[]string{id},
		"",
		sc.DeleteSentence,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBからレコードが削除されている
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences
		WHERE id = $1;
	`,
		id,
	).Scan(&count)

	assert.Equal(t, 0, count)
}

func TestDeleteSentenceWithoutLoggingInUserId(t *testing.T) {
	// ログイン中のUserに紐づかないSentenceを削除できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ削除可能とする
	DeleteAllFromSentences()

	var id string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 2)
		RETURNING id;
	`).Scan(&id)

	DoSimpleTest(
		t,
		http.MethodDelete,
		"/sentences/:sentenceId",
		[]string{"sentenceId"},
		[]string{id},
		"",
		sc.DeleteSentence,
		http.StatusUnauthorized,
		"{}",
	)

	// DBからレコードが削除されていない
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences
		WHERE id = $1;
	`,
		id,
	).Scan(&count)

	assert.Equal(t, 1, count)
}
func TestAssociateSentenceWithWords(t *testing.T) {
	// WordとSentenceがどちらもログイン中のUserに紐づく場合
	// それらを紐づかせられることをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo',  1)
		RETURNING id;
	`).Scan(&wordId)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := fmt.Sprintf(`
		{
			"word_ids": [%s]
		}
	`,
		wordId)

	expectedResponse := fmt.Sprintf(`
		{
			"word_ids": [%s]
		}`,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/association/:sentenceId",
		[]string{"sentenceId"},
		[]string{sentenceId},
		reqBody,
		sc.AssociateSentenceWithWords,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBにレコードが追加されている
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		wordId,
	).Scan(&count)

	assert.Equal(t, 1, count)
}

func TestAssociateSentenceWithWordsWithMultipleWordIds(t *testing.T) {
	// WordとSentenceがどちらもログイン中のUserに紐づく場合
	// かつWordが複数選択されている場合
	// それらを紐づかせられることをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId1 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word1', 'memo1',  1)
		RETURNING id;
	`).Scan(&wordId1)

	var wordId2 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word2', 'memo2',  1)
		RETURNING id;
	`).Scan(&wordId2)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := fmt.Sprintf(`
		{
			"word_ids": [%s, %s]
		}`,
		wordId1, wordId2,
	)

	expectedResponse := fmt.Sprintf(`
		{
			"word_ids": [%s, %s]
		}`,
		wordId1, wordId2,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/association/:sentenceId",
		[]string{"sentenceId"},
		[]string{sentenceId},
		reqBody,
		sc.AssociateSentenceWithWords,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBにレコードが追加されている
	var countWithWordId1 int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		wordId1,
	).Scan(&countWithWordId1)

	assert.Equal(t, 1, countWithWordId1)

	var countWithWordId2 int
	db.QueryRow(`
	SELECT COUNT(*) FROM sentences_words
	WHERE sentence_id = $1
		AND word_id = $2;
	`,
		sentenceId,
		wordId2,
	).Scan(&countWithWordId2)

	assert.Equal(t, 1, countWithWordId2)
}

func TestAssociateSentenceWithInvalidWordId(t *testing.T) {
	// Sentenceはログイン中のUserに紐づくが
	// Wordがログイン中のUserに紐づかない場合
	// それらを紐づかせられないことをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo',  2)
		RETURNING id;
	`).Scan(&wordId)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := fmt.Sprintf(`
		{
			"word_ids": [%s]
		}`,
		wordId,
	)

	expectedResponse := `{
		"word_ids": null
	}`

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/association/:sentenceId",
		[]string{"sentenceId"},
		[]string{sentenceId},
		reqBody,
		sc.AssociateSentenceWithWords,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBにレコードが追加されていない
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		wordId,
	).Scan(&count)

	assert.Equal(t, 0, count)
}

func TestAssociateSentenceWithInvalidSentenceId(t *testing.T) {
	// Wordはログイン中のUserに紐づくが
	// Sentenceがログイン中のUserに紐づかない場合
	// それらを紐づかせられないことをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word', 'memo',  1)
		RETURNING id;
	`).Scan(&wordId)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 2)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := fmt.Sprintf(`
		{
			"word_ids": [%s]
		}`,
		wordId,
	)

	expectedResponse := `{
		"word_ids": null
	}`

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/association/:sentenceId",
		[]string{"sentenceId"},
		[]string{sentenceId},
		reqBody,
		sc.AssociateSentenceWithWords,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBにレコードが追加されていない
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		wordId,
	).Scan(&count)

	assert.Equal(t, 0, count)
}

func TestAssociateSentenceWithAllInvalidWordId(t *testing.T) {
	// Sentenceはログイン中のUserに紐づくが
	// 指定された複数のWordのうち、すべてがログイン中のUserに紐づかない場合
	// それらを紐づかせられないことをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var invalidWordId1 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word1', 'memo1',  2)
		RETURNING id;
	`).Scan(&invalidWordId1)

	var invalidWordId2 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word2', 'memo2',  2)
		RETURNING id;
	`).Scan(&invalidWordId2)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := fmt.Sprintf(`
		{
			"word_ids": [%s, %s]
		}`,
		invalidWordId1,
		invalidWordId2,
	)

	expectedResponse := `{
		"word_ids": null
	}`

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/association/:sentenceId",
		[]string{"sentenceId"},
		[]string{sentenceId},
		reqBody,
		sc.AssociateSentenceWithWords,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBにレコードが追加されていない
	var countWithInvalidWordId1 int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		invalidWordId1,
	).Scan(&countWithInvalidWordId1)

	assert.Equal(t, 0, countWithInvalidWordId1)

	var countWithInvalidWordId2 int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		invalidWordId2,
	).Scan(&countWithInvalidWordId2)

	assert.Equal(t, 0, countWithInvalidWordId2)
}

func TestAssociateSentenceWithSomeInvalidWordId(t *testing.T) {
	// Sentenceはログイン中のUserに紐づくが
	// 指定された複数のWordのうち、一部だけがログイン中のUserに紐づかない場合
	// それらを紐づかせられないことをテスト
	// TODO ログイン機能
	// とりあえずWordとSentenceのUserIdが両方1の場合紐づけ可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var validWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word1', 'memo1', 1)
		RETURNING id;
	`).Scan(&validWordId)

	var invalidWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'word2', 'memo2',  2)
		RETURNING id;
	`).Scan(&invalidWordId)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('word_id_seq'), 'sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := fmt.Sprintf(`
		{
			"word_ids": [%s, %s]
		}`,
		validWordId,
		invalidWordId,
	)

	expectedResponse := fmt.Sprintf(`
		{
			"word_ids": [%s]
		}`,
		validWordId,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/association/:sentenceId",
		[]string{"sentenceId"},
		[]string{sentenceId},
		reqBody,
		sc.AssociateSentenceWithWords,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBにレコードが追加されていない
	var countWithValidWordId int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		validWordId,
	).Scan(&countWithValidWordId)

	assert.Equal(t, 1, countWithValidWordId)

	var countWithInvalidWordId int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
	`,
		sentenceId,
		invalidWordId,
	).Scan(&countWithInvalidWordId)

	assert.Equal(t, 0, countWithInvalidWordId)
}

func TestGetAssociatedWords(t *testing.T) {
	// WordとSentenceがどちらもログイン中のuser_idに紐づく場合、
	// Sentenceに紐づくWordを取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	var wordId1 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word1', 'test memo1', 1)
		RETURNING id;
	`).Scan(&wordId1)

	var wordId2 string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word2', 'test memo2', 1)
		RETURNING id;
	`).Scan(&wordId2)

	db.QueryRow(`
		INSERT INTO sentences_words
		(sentence_id, word_id)
		VALUES
		($1, $2),
		($1, $3);
		`,
		sentenceId,
		wordId1,
		wordId2,
	)

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %s,
				"word": "test word1",
				"memo": "test memo1",
				"user_id": 1
			},
			{
				"id": %s,
				"word": "test word2",
				"memo": "test memo2",
				"user_id": 1
			}
		]`,
		wordId1,
		wordId2,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences/:sentenceId/associated-words",
		[]string{"sentenceId"},
		[]string{sentenceId},
		"",
		sc.GetAssociatedWords,
		http.StatusOK,
		expectedResponse,
	)
}

func TestGetAssociatedWordsWithInvalidSentenceId(t *testing.T) {
	// Sentenceがログイン中のuser_idに紐づかない場合、
	// Sentenceに紐づくWordを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	// Sentenceのuser_idとしてログイン中のもの以外を使用
	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence', 2)
		RETURNING id;
	`).Scan(&sentenceId)

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word1', 'test memo1', 1)
		RETURNING id;
	`).Scan(&wordId)

	db.QueryRow(`
		INSERT INTO sentences_words
		(sentence_id, word_id)
		VALUES
		($1, $2);
		`,
		sentenceId,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences/:sentenceId/associated-words",
		[]string{"sentenceId"},
		[]string{sentenceId},
		"",
		sc.GetAssociatedWords,
		http.StatusOK,
		"null",
	)
}

func TestGetAssociatedWordsWithInvalidWordId(t *testing.T) {
	// Wordがログイン中のuser_idに紐づかない場合、
	// Sentenceに紐づくWordを取得できないことをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	// Wordのuser_idとしてログイン中のもの以外を使用
	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'test word', 'test memo', 2)
		RETURNING id;
	`).Scan(&wordId)

	db.QueryRow(`
		INSERT INTO sentences_words
		(sentence_id, word_id)
		VALUES
		($1, $2);
		`,
		sentenceId,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/sentences/:sentenceId/associated-words",
		[]string{"sentenceId"},
		[]string{sentenceId},
		"",
		sc.GetAssociatedWords,
		http.StatusOK,
		"null",
	)
}
