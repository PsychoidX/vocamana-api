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
		"/sentences",
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
		VALUES(nextval('sentence_id_seq'), 'test sentence', 1)
		RETURNING id;
	`).Scan(&idWithUserId1)

	var idWithUserId2 int
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'test sentence', 2)
		RETURNING id;
	`).Scan(&idWithUserId2)

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"sentence": "test sentence",
				"sentence_with_link": "test sentence",
				"user_id": 1
			}
		]`,
		idWithUserId1,
	)

	DoSimpleTest(
		t,
		"/sentences",
		sc.GetAllSentences,
		http.StatusOK,
		expectedResponse,
	)
}

func TestGetAllSentences_IncludingWords(t *testing.T) {
	// 取得したSentenceに正しくリンクが含まれていることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	appleWordId := createTestWord(t, "林檎", "").Id
	ateWordId := createTestWord(t, "食べた", "").Id
	sentenceId := createTestSentence(t, "林檎を食べた").Id

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"sentence": "林檎を食べた",
				"sentence_with_link": "<a href=\"/words/%d\">林檎</a>を<a href=\"/words/%d\">食べた</a>",
				"user_id": 1
			}
		]`,
		sentenceId,
		appleWordId,
		ateWordId,
	)

	DoSimpleTest(
		t,
		"/sentences",
		sc.GetAllSentences,
		http.StatusOK,
		expectedResponse,
	)
}

func TestGetAllSentences_WithLimit(t *testing.T) {
	// ログイン中のUserに紐づくSentenceを、LIMIT付きで取得できることをテテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromSentences()

	sentenceId1 := createTestSentence(t, "test sentence 1").Id
	sentenceId2 := createTestSentence(t, "test sentence 2").Id
	createTestSentence(t, "test sentence 3")

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"sentence": "test sentence 1",
				"sentence_with_link": "test sentence 1",
				"user_id": 1
			},
			{
				"id": %d,
				"sentence": "test sentence 2",
				"sentence_with_link": "test sentence 2",
				"user_id": 1
			}
		]`,
		sentenceId1,
		sentenceId2,
	)

	DoSimpleTest(
		t,
		"/sentences",
		sc.GetAllSentences,
		http.StatusOK,
		expectedResponse,
		QueryParams(
			[]string{"limit"},
			[][]string{{"2"}},
		),
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
		"/sentences/:sentenceId",
		sc.GetSentenceById,
		http.StatusOK,
		expectedResponse,
		Params(
			[]string{"sentenceId"},
			[]string{id},
		),
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
		"/sentences/:sentenceId",
		sc.GetSentenceById,
		http.StatusOK,
		"{}",
		Params(
			[]string{"sentenceId"},
			[]string{id},
		),
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
		"/sentences",
		sc.CreateSentence,
		http.StatusCreated,
		expectedResponse,
		Body(reqBody),
		HttpMethod(http.MethodPost),
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
		"/sentences",
		sc.CreateSentence,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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
		"/sentences",
		sc.CreateSentence,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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
		"/sentences",
		sc.CreateSentence,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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
		"/sentences",
		sc.CreateSentence,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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
		"/sentences/multiple",
		sc.CreateMultipleSentences,
		http.StatusCreated,
		expectedResponse,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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
		"/words/:sentenceId",
		sc.UpdateSentence,
		http.StatusAccepted,
		expectedResponse,
		HttpMethod(http.MethodPut),
		Params(
			[]string{"sentenceId"},
			[]string{id},
		),
		Body(reqBody),
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

	redWordId := createTestWord(t, "赤い", "").Id
	blueWordId := createTestWord(t, "青い", "").Id
	sentenceId := createTestSentence(t, "赤いりんごを食べた").Id

	updateSentenceReqBody := `{
		"sentence": "青いりんごを食べた"
	}`

	ExecController(
		t,
		"/words/:sentenceId",
		sc.UpdateSentence,
		Params(
			[]string{"sentenceId"},
			[]string{strconv.FormatUint(sentenceId, 10)},
		),
		Body(updateSentenceReqBody),
		HttpMethod(http.MethodPut),
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
		"/words/:sentenceId",
		sc.UpdateSentence,
		http.StatusUnauthorized,
		"{}",
		Params(
			[]string{"sentenceId"},
			[]string{id},
		),	
		Body(reqBody),
		HttpMethod(http.MethodPut),
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
		"/sentences/:sentenceId",
		sc.DeleteSentence,
		http.StatusAccepted,
		expectedResponse,
		Params(
			[]string{"sentenceId"},
			[]string{id},
		),
		HttpMethod(http.MethodDelete),
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

	wordId := createTestWord(t, "赤い", "").Id
	sentenceId := createTestSentence(t, "赤いりんごを食べた").Id

	ExecController(
		t,
		"/words/:sentenceId",
		sc.DeleteSentence,
		Params(
			[]string{"sentenceId"},
			[]string{strconv.FormatUint(sentenceId, 10)},
		),
		HttpMethod(http.MethodDelete),
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
		"/sentences/:sentenceId",
		sc.DeleteSentence,
		http.StatusUnauthorized,
		"{}",
		Params(
			[]string{"sentenceId"},
			[]string{id},
		),
		HttpMethod(http.MethodDelete),
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
		"/sentences/:sentenceId/associated-words",
		sc.GetAssociatedWords,
		http.StatusOK,
		expectedResponse,
		Params(
			[]string{"sentenceId"},
			[]string{sentenceId},
		),
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
		"/sentences/:sentenceId/associated-words",
		sc.GetAssociatedWords,
		http.StatusOK,
		"null",
		Params(
			[]string{"sentenceId"},
			[]string{sentenceId},
		),
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
		"/sentences/:sentenceId/associated-words",
		sc.GetAssociatedWords,
		http.StatusOK,
		"null",
		Params(
			[]string{"sentenceId"},
			[]string{sentenceId},
		),
	)
}

func TestGetSentencesCount(t *testing.T) {
	// ログイン中のユーザの単語数を取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromSentences()

	DoSimpleTest(
		t,
		"/sentences/count",
		sc.GetSentencesCount,
		http.StatusOK,
		`{ "count": 0 }`,
	)

	createTestSentence(t, "test sentence")

	DoSimpleTest(
		t,
		"/sentences/count",
		sc.GetSentencesCount,
		http.StatusOK,
		`{ "count": 1 }`,
	)

	createTestSentence(t, "test sentence")
	createTestSentence(t, "test sentence")

	DoSimpleTest(
		t,
		"/sentences/count",
		sc.GetSentencesCount,
		http.StatusOK,
		`{ "count": 3 }`,
	)
}