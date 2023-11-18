package test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllSentences_WithNoRows(t *testing.T) {
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

func TestGetSentenceById(t *testing.T) {
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

func TestGetSentenceById_WithInvalidUser(t *testing.T) {
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

func TestCreateSentence_IncludingWords(t *testing.T) {
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

	var ateWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), '食べた', 'test memo', 1)
		RETURNING id;
	`).Scan(&ateWordId)

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

	// 該当するWordがsentences_wordsに追加されることをテスト：
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

	// 該当するNotationがsentences_wordsに追加されることをテスト：
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

	// 該当しないWordがsentences_wordsに追加されないことをテスト：
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

	// 2つ目の該当するWordがsentences_wordsに追加されることをテスト：
	// 「赤いりんごを食べた」には「食べた」が含まれるため、
	// sentences_wordsに追加される
	var ateCount int
	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		ateWordId,
	).Scan(&ateCount)

	assert.Equal(t, 1, ateCount)
}

func TestCreateSentence_IncludingInvalidWords(t *testing.T) {
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

func TestCreateSentence_IncludingNotations(t *testing.T) {
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

func TestCreateSentence_IncludingInvalidNotations(t *testing.T) {
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

func TestCreateMultipleSentences(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを複数同時に作成できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromSentences()

	sentenceId1 := GetNextSentencesSequenceValue()
	sentenceId2 := sentenceId1 + 1

	reqBody := `{
		"sentences": [
				{
					"sentence": "test sentence 1"
				},
				{
					"sentence": "test sentence 2"
				}
			]
		}`

	// 登録されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"sentence": "test sentence 1",
				"user_id": 1
			},
			{
				"id": %d,
				"sentence": "test sentence 2",
				"user_id": 1
			}
		]`,
		sentenceId1,
		sentenceId2,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/sentences/multiple",
		nil,
		nil,
		reqBody,
		sc.CreateMultipleSentences,
		http.StatusCreated,
		expectedResponse,
	)

	// DBにレコードが追加される
	var sentence1 string
	db.QueryRow(`
		SELECT sentence FROM sentences
		WHERE id = $1;
	`,
		sentenceId1,
	).Scan(&sentence1)

	assert.Equal(t, "test sentence 1", sentence1)

	var sentence2 string
	db.QueryRow(`
		SELECT sentence FROM sentences
		WHERE id = $1;
	`,
		sentenceId2,
	).Scan(&sentence2)

	assert.Equal(t, "test sentence 2", sentence2)
}

func TestUpdateSentence(t *testing.T) {
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

func TestUpdateSentence_UpdatedAssociation(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを更新した時、
	// sentences_wordsが正常に再構築されることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ更新可能とする
	DeleteAllFromSentences()
	DeleteAllFromWords()

	redWordRes := createTestWord(t, "赤い", "")
	redWordId := redWordRes.Id

	blueWordRes := createTestWord(t, "青い", "")
	blueWordId := blueWordRes.Id

	sentenceRes := createTestSentence(t, "赤いりんごを食べた")
	sentenceId := sentenceRes.Id

	updateSentenceReqBody := `{
		"sentence": "青いりんごを食べた"
	}`
	
	ExecController(
		t,
		http.MethodPut,
		"/words/:sentenceId",
		[]string{"sentenceId"},
		[]string{strconv.FormatUint(sentenceId, 10)},
		updateSentenceReqBody,
		sc.UpdateSentence,
	)

	// Sentenceを変更したことで、
	// 単語「赤い」がSentence中に含まれなくなるため
	// sentences_wordsからも削除される
	assert.Equal(t, 0, getCountFromSentencesWords(sentenceId, redWordId))

	// 単語「青い」がSentence中に含まれるようになるため
	// sentences_wordsに追加される
	assert.Equal(t, 1, getCountFromSentencesWords(sentenceId, blueWordId))
}

func TestUpdateSentence_WithInvalidUser(t *testing.T) {
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

func TestDeleteSentence(t *testing.T) {
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

func TestDeleteSentence_UpdateAssociation(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを削除した時、
	// sentences_wordsから、削除されたSentenceに関するレコードが正常に削除されることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ削除可能とする
	DeleteAllFromSentences()
	DeleteAllFromWords()

	wordRes := createTestWord(t, "赤い", "")
	wordId := wordRes.Id

	sentenceRes := createTestSentence(t, "赤いりんごを食べた")
	sentenceId := sentenceRes.Id
	
	ExecController(
		t,
		http.MethodDelete,
		"/words/:sentenceId",
		[]string{"sentenceId"},
		[]string{strconv.FormatUint(sentenceId, 10)},
		"",
		sc.DeleteSentence,
	)

	// Sentenceを変更した時、
	// sentences_wordsからも削除される
	assert.Equal(t, 0, getCountFromSentencesWords(sentenceId, wordId))
}

func TestDeleteSentence_WithInvalidUser(t *testing.T) {
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

func TestGetAssociatedWords_WithInvalidSentenceId(t *testing.T) {
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

func TestGetAssociatedWords_WithInvalidWordId(t *testing.T) {
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