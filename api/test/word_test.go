package test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllWords_WithNoRows(t *testing.T) {
	// ログイン中のUserに紐づくWordが1つも無い場合nullが返ることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のWordのみ取得可能とする
	DeleteAllFromWords()

	// レコードが1つも無い場合、[]ではなくnullが返る
	DoSimpleTest(
		t,
		"/words",
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
		"/words",
		wc.GetAllWords,
		http.StatusOK,
		expectedResponse,
	)
}

func TestGetWordById(t *testing.T) {
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
		"/words/:wordId",
		wc.GetWordById,
		http.StatusOK,
		expectedResponse,
		Params(
			[]string{"wordId"},
			[]string{id},
		),
	)
}

func TestGetWordById_WithInvalidUser(t *testing.T) {
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
		"/words/:wordId",
		wc.GetWordById,
		http.StatusOK,
		"{}",
		Params(
			[]string{"wordId"},
			[]string{id},
		),
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
		"/words",
		wc.CreateWord,
		http.StatusCreated,
		expectedResponse,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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
		"/words/multiple",
		wc.CreateMultipleWords,
		http.StatusCreated,
		expectedResponse,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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

func TestCreateWord_InSentences(t *testing.T) {
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
		"/words",
		wc.CreateWord,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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

func TestCreateWord_InInvalidSentences(t *testing.T) {
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
		"/words",
		wc.CreateWord,
		HttpMethod(http.MethodPost),
		Body(reqBody),
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

func TestCreateWord_RootWordCreated(t *testing.T) {
	// 新規Wordを追加したとき、語幹のNotationも追加されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	// 「買う」の語幹「買」が追加される
	wordId := createTestWord(t, "買う", "").Id
	assert.Equal(t, 1, getCountFromNotationsByNotation(wordId, "買"))
}

func TestUpdateWord(t *testing.T) {
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
		"/words/:wordId",
		wc.UpdateWord,
		http.StatusAccepted,
		expectedResponse,
		Params(
			[]string{"wordId"},
			[]string{id},
		),
		HttpMethod(http.MethodPut),
		Body(reqBody),
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

func TestUpdateWord_UpdatedAssociation(t *testing.T) {
	// ログイン中のUserに紐づくWordを更新した時、
	// sentences_wordsが正常に再構築されることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ更新可能とする
	DeleteAllFromSentences()
	DeleteAllFromWords()

	wordId := createTestWord(t, "赤い", "").Id
	appleSentenceId := createTestSentence(t, "赤いリンゴを食べた").Id
	lemonSentenceId := createTestSentence(t, "黄色いレモンを食べた").Id

	reqBody := `{
		"word": "黄色い",
		"memo": ""
	}`

	ExecController(
		t,
		"/words/:wordeId",
		wc.UpdateWord,
		HttpMethod(http.MethodPut),
		Params(
			[]string{"wordId"},
			[]string{strconv.FormatUint(wordId, 10)},
		),
		Body(reqBody),
	)

	// Wordを変更したことで、
	// 「赤いリンゴを食べた」の中にWordが含まれなくなるため
	// sentences_wordsから削除される
	assert.Equal(t, 0, getCountFromSentencesWords(appleSentenceId, wordId))

	// 「黄色いレモンを食べた」の中にWordが含まれるようになるため
	// sentences_wordsに追加される
	assert.Equal(t, 1, getCountFromSentencesWords(lemonSentenceId, wordId))
}

func TestUpdateWord_WithInvalidUser(t *testing.T) {
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
		"/words/:wordId",
		wc.UpdateWord,
		http.StatusUnauthorized,
		"{}",
		HttpMethod(http.MethodPut),
		Params(
			[]string{"wordId"},
			[]string{id},
		),
		Body(reqBody),
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

func TestUpdateWord_RootWordCreated(t *testing.T) {
	// Wordを更新したとき、語幹のNotationも更新されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	// 「買う」の語幹「買」が追加される
	wordId := createTestWord(t, "買う", "").Id
	assert.Equal(t, 1, getCountFromNotationsByNotation(wordId, "買"))

	reqBody := `{
		"word": "赤い",
		"memo": ""
	}`

	ExecController(
		t,
		"/words/:wordeId",
		wc.UpdateWord,
		Params(
			[]string{"wordId"},
			[]string{strconv.FormatUint(wordId, 10)},
		),
		Body(reqBody),
		HttpMethod(http.MethodPut),
	)

	// 「買う」の語幹「買」が削除される
	assert.Equal(t, 0, getCountFromNotationsByNotation(wordId, "買"))
	// 「赤い」の語幹「赤」が追加される
	assert.Equal(t, 1, getCountFromNotationsByNotation(wordId, "赤"))
}

func TestUpdateWord_RootWordCreated_OldRootWordNotExists(t *testing.T) {
	// Wordを更新したとき、更新前のWordの語幹のNotationが存在しなくても、正常に更新されることをテスト

	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ作成可能とする
	DeleteAllFromWords()
	DeleteAllFromSentences()

	wordId := insertIntoWords("買う", "", 1)
	assert.Equal(t, 0, getCountFromNotationsByNotation(wordId, "買"))

	reqBody := `{
		"word": "赤い",
		"memo": ""
	}`

	ExecController(
		t,
		"/words/:wordeId",
		wc.UpdateWord,
		Params(
			[]string{"wordId"},
			[]string{strconv.FormatUint(wordId, 10)},
		),
		Body(reqBody),
		HttpMethod(http.MethodPut),
	)

	// 「赤い」の語幹「赤」が追加される
	assert.Equal(t, 1, getCountFromNotationsByNotation(wordId, "赤"))
}

func TestDeleteWord(t *testing.T) {
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
		"/words/:wordId",
		wc.DeleteWord,
		http.StatusAccepted,
		expectedResponse,
		Params(
			[]string{"wordId"},
			[]string{id},
		),
		HttpMethod(http.MethodDelete),
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

func TestDeleteWord_WithInvalidUser(t *testing.T) {
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
		"/words/:wordId",
		wc.DeleteWord,
		http.StatusUnauthorized,
		"{}",
		HttpMethod(http.MethodDelete),
		Params(
			[]string{"wordId"},
			[]string{id},
		),
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

func TestDeleteWord_UpdatedAssociation(t *testing.T) {
	// ログイン中のUserに紐づくWordを削除した時、
	// sentences_wordsが正常に再構築されることをテスト
	// TODO ログイン機能
	// とりあえずuser_id=1のSentenceのみ更新可能とする
	DeleteAllFromSentences()
	DeleteAllFromWords()

	wordId := createTestWord(t, "赤い", "").Id
	sentenceId := createTestSentence(t, "赤いリンゴを食べた").Id

	ExecController(
		t,
		"/words/:wordeId",
		wc.DeleteWord,
		Params(
			[]string{"wordId"},
			[]string{strconv.FormatUint(wordId, 10)},
		),
		HttpMethod(http.MethodDelete),
	)

	// Wordを削除したことで、
	// 「赤いリンゴを食べた」の中にWordが含まれなくなるため
	// sentences_wordsから削除される
	assert.Equal(t, 0, getCountFromSentencesWords(sentenceId, wordId))
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
		"/words/:wordId/associated-sentences",
		wc.GetAssociatedSentencesWithLink,
		http.StatusOK,
		expectedResponse,
		Params(
			[]string{"wordId"},
			[]string{appleWordId},
		),
	)
}