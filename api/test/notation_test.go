package test

import (
	"fmt"
	"net/http"
	"strconv"
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

func TestGetAllNotations_WithNoRows(t *testing.T) {
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

func TestGetAllNotations_WithInvalidWordId(t *testing.T) {
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

func TestCreateNotation(t *testing.T) {
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

func TestCreateNotation_WithInvalidUser(t *testing.T) {
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

func TestCreateNotation_InSentence(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを追加をした時、
	// 既存のSentence中に新規Notationを含むものがある場合、
	// sentences_wordsに追加されることをテスト

	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromSentences()
	DeleteAllFromNotations()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'りんご', 'test memo', 1)
		RETURNING id;
	`).Scan(&wordId)

	var appleSentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), '赤い林檎を食べた', 1)
		RETURNING id;
	`).Scan(&appleSentenceId)

	var lemonSentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), '黄色い檸檬を食べた', 1)
		RETURNING id;
	`).Scan(&lemonSentenceId)

	reqBody := `{
		"notation": "林檎"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordId},
		reqBody,
		nc.CreateNotation,
	)

	// 「赤い林檎を食べた」には「林檎」が含まれるため、
	// sentences_wordsに追加される
	var appleCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		appleSentenceId,
		wordId,
	).Scan(&appleCount)

	assert.Equal(t, 1, appleCount)

	// 「黄色い檸檬を食べた」には「林檎」が含まれないため、
	// sentences_wordsに追加されない
	var lemonCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		lemonSentenceId,
		wordId,
	).Scan(&lemonCount)

	assert.Equal(t, 0, lemonCount)
}

func TestCreateNotation_InInvalidSentence(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを追加をした時、
	// 既存のSentence中に新規Notationを含むものがあるが、
	// 該当Sentenceのuser_idがログイン中のものと異なる場合、
	// sentences_wordsに追加されることをテスト

	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromSentences()
	DeleteAllFromNotations()

	var wordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'りんご', 'test memo', 1)
		RETURNING id;
	`).Scan(&wordId)

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), '赤い林檎を食べた', 2)
		RETURNING id;
	`).Scan(&sentenceId)

	reqBody := `{
		"notation": "林檎"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{wordId},
		reqBody,
		nc.CreateNotation,
	)

	// 「赤い林檎を食べた」には「林檎」が含まれるが、
	// user_idが異なるため、
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
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{notationId},
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

func TestUpdateNotation_WithNoRows(t *testing.T) {
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
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{"1"},
		reqBody,
		nc.UpdateNotation,
		http.StatusUnauthorized,
		"{}",
	)
}

func TestDeleteNotation(t *testing.T) {
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
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{notationId},
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

func TestDeleteNotation_WithInvalidUser(t *testing.T) {
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
		"/notations/:notationId",
		[]string{"notationId"},
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

func TestUpdateNotation_UpdatedAssociation(t *testing.T) {
	// Notationを更新した時、entences_wordsが正常に再構築されることをテスト
	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := createTestWord(t, "りんご", "").Id
	notationId := createTestNotation(t, wordId, "林檎").Id
	sentenceId := createTestSentence(t, "林檎を食べた").Id

	body := `{
		"notation": "リンゴ"
	}`

	ExecController(
		t,
		http.MethodPut,
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{strconv.FormatUint(notationId, 10)},
		body,
		nc.UpdateNotation,
	)

	// Word「りんご」が、追加されていたNotationにより「林檎」にもマッチしていたが、
	// Notationを「林檎」から「リンゴ」に変更したことで、「林檎」にはマッチしなくなり
	// sentences_wordsから削除される
	assert.Equal(t, 0, getCountFromSentencesWords(sentenceId, wordId))
}

func TestDeleteNotation_UpdatedAssociation_AllNotationsDeleted(t *testing.T) {
	// Notationを削除した時、それによってSentence中にWordが含まれないことになったら、
	// sentences_wordsから値が削除されることをテスト

	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := createTestWord(t, "りんご", "").Id
	notationId := createTestNotation(t, wordId, "林檎").Id
	sentenceId := createTestSentence(t, "林檎を食べた").Id

	ExecController(
		t,
		http.MethodDelete,
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{strconv.FormatUint(notationId, 10)},
		"",
		nc.DeleteNotation,
	)

	// Word「りんご」が、追加されていたNotationにより「林檎」にもマッチしていたが、
	// Notation「林檎」を削除したことで、Wordが「林檎」にマッチしなくなり、
	// sentences_wordsから削除される
	assert.Equal(t, 0, getCountFromSentencesWords(sentenceId, wordId))
}

func TestDeleteNotation_UpdatedAssociation_SomeNotationRemaining(t *testing.T) {
	// Notationを削除した時、それでもSentence中にWordが含まれている場合、
	// sentences_wordsから値が削除されないことをテスト

	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := createTestWord(t, "りんご", "").Id
	notationId1 := createTestNotation(t, wordId, "林檎").Id
	createTestNotation(t, wordId, "リンゴ")
	sentenceId := createTestSentence(t, "林檎はリンゴと読むらしい").Id

	ExecController(
		t,
		http.MethodDelete,
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{strconv.FormatUint(notationId1, 10)},
		"",
		nc.DeleteNotation,
	)

	// Word「りんご」が、追加されていたNotationにより「林檎」「リンゴ」にもマッチしているので、
	// Notation「林檎」を削除しても、Wordは「リンゴ」にマッチし続けるため、
	// sentences_wordsから削除されない
	assert.Equal(t, 1, getCountFromSentencesWords(sentenceId, wordId))
}