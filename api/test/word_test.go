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

func TestCreateMultipleWords(t *testing.T) {
	// ログイン中のUserに紐づくWordを複数同時に作成できることをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()

	wordId1 := GetNextWordsSequenceValue()
	wordId2 := wordId1 + 1

	reqBody := `{
		"words": [
			{
				"word": "test word 1",
				"memo": "test memo 1"
			},
			{
				"word": "test word 2",
				"memo": "test memo 2"
			}
		]
	}`

	// 登録されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"word": "test word 1",
				"memo": "test memo 1",
				"user_id": 1
			},
			{
				"id": %d,
				"word": "test word 2",
				"memo": "test memo 2",
				"user_id": 1
			}
		]`,
		wordId1,
		wordId2,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/words/multiple",
		nil,
		nil,
		reqBody,
		wc.CreateMultipleWords,
		http.StatusCreated,
		expectedResponse,
	)

	// DBにレコードが追加される
	var word1 string
	var memo1 string
	db.QueryRow(`
		SELECT word, memo FROM words
		WHERE id = $1;
		`,
		wordId1,
	).Scan(&word1, &memo1)

	assert.Equal(t, "test word 1", word1)
	assert.Equal(t, "test memo 1", memo1)

	var word2 string
	var memo2 string
	db.QueryRow(`
		SELECT word, memo FROM words
		WHERE id = $1;
		`,
		wordId2,
	).Scan(&word2, &memo2)

	assert.Equal(t, "test word 2", word2)
	assert.Equal(t, "test memo 2", memo2)
}

func TestCreateWordInSentences(t *testing.T) {
	// 既存のSentence中に、新規追加したWordを含むものがある場合、
	// sentences_wordsに追加されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var appleSentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), '赤いりんごを食べた', 1)
		RETURNING id;
	`).Scan(&appleSentenceId)

	var lemonSentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), '黄色いレモンを食べた', 1)
		RETURNING id;
	`).Scan(&lemonSentenceId)

	wordId := GetNextWordsSequenceValue()

	reqBody := `{
		"word": "りんご",
		"memo": "test memo"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/words",
		nil,
		nil,
		reqBody,
		wc.CreateWord,
	)

	// 「赤いりんごを食べた」には「赤い」が含まれるため、
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

	// 「黄色いレモンを食べた」には「りんご」が含まれないため、
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

func TestCreateWordInInvalidSentences(t *testing.T) {
	// 既存のSentence中に、新規追加したWordを含むものがあるが、
	// ログイン中のuser_idと異なる場合、
	// sentences_wordsに追加されないことをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), '赤いりんごを食べた', 2)
		RETURNING id;
	`).Scan(&sentenceId)

	wordId := GetNextWordsSequenceValue()

	reqBody := `{
		"word": "りんご",
		"memo": "test memo"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/words",
		nil,
		nil,
		reqBody,
		wc.CreateWord,
	)

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
	// Sentenceを、WordとNotationがaタグに変換された状態で取得できることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ取得可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	var sentenceId string
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), 'リンゴと林檎、レモンと檸檬が同一であるとみなす', 1)
		RETURNING id;
	`).Scan(&sentenceId)

	var appleWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'リンゴ', 'word memo', 1)
		RETURNING id;
	`).Scan(&appleWordId)

	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES
		(nextval('word_id_seq'), $1, '林檎');
		`,
		appleWordId,
	)

	var lemonWordId string
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), 'レモン', 'word memo', 1)
		RETURNING id;
	`).Scan(&lemonWordId)

	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES
		(nextval('word_id_seq'), $1, '檸檬');
		`,
		lemonWordId,
	)

	db.QueryRow(`
		INSERT INTO sentences_words
		(sentence_id, word_id)
		VALUES
		($1, $2),
		($1, $3);
		`,
		sentenceId,
		appleWordId,
		lemonWordId,
	)

	// :wordId == appleWordId では、紐づくSentenceとして「リンゴと林檎、レモンと檸檬が同一であるとみなす」が取得される
	// このうち「レモン」「檸檬」のwordIdはappleWordIdと同値でないが、
	// APIアクセスに使用したwordIdでない単語（レモン・檸檬）に関しても、userIdが同じであれば、リンク作成はされる
	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %s,
				"sentence": "リンゴと林檎、レモンと檸檬が同一であるとみなす",
				"sentence_with_link": "<a href=\"/words/%s\">リンゴ</a>と<a href=\"/words/%s\">林檎</a>、<a href=\"/words/%s\">レモン</a>と<a href=\"/words/%s\">檸檬</a>が同一であるとみなす",
				"user_id": 1
			}
		]`,
		sentenceId,
		appleWordId,
		appleWordId,
		lemonWordId,
		lemonWordId,
	)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/associated-sentences",
		[]string{"wordId"},
		[]string{appleWordId},
		"",
		wc.GetAssociatedSentencesWithLink,
		http.StatusOK,
		expectedResponse,
	)
}
