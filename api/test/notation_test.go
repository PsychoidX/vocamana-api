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

	wordId := insertIntoWords("test word", "test memo", 1)
	notationId1 := insertIntoNotations(wordId, "test notation1")
	notationId2 := insertIntoNotations(wordId, "test notation2")

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"word_id": %d,
				"notation": "test notation1"
			},
			{
				"id": %d,
				"word_id": %d,
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
		[]string{strconv.FormatUint(wordId, 10)},
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

	wordId := insertIntoWords("test word", "test memo", 1)

	DoSimpleTest(
		t,
		http.MethodGet,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{strconv.FormatUint(wordId, 10)},
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

	wordIdWithUserId1 := insertIntoWords("test word", "test memo", 1)
	wordIdWithUserId2 := insertIntoWords("test word", "test memo", 2)
	notationIdWithUserId1 := insertIntoNotations(wordIdWithUserId1, "test notation1")
	insertIntoNotations(wordIdWithUserId2, "test notation2")

	expectedResponse := fmt.Sprintf(`
		[
			{
				"id": %d,
				"word_id": %d,
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
		[]string{strconv.FormatUint(wordIdWithUserId1, 10)},
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
		[]string{strconv.FormatUint(wordIdWithUserId2, 10)},
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

	wordId := insertIntoWords("test word", "test memo", 1)

	notationId := GetNextNotationsSequenceValue()

	reqBody := `{
		"notation": "test notation"
	}`

	// 登録されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %d,
			"word_id": %d,
			"notation": "test notation"
		}`,
		notationId,
		wordId,
	)

	DoSimpleTest(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{strconv.FormatUint(wordId, 10)},
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
		notationId,
		wordId,
	).Scan(&notation)

	assert.Equal(t, "test notation", notation)
}

func TestCreateNotation_WithInvalidUser(t *testing.T) {
	// ログイン中のUserに紐づかないWordに対し、Notationを作成できないことをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := insertIntoWords("test word", "test memo", 2)

	notationId := GetNextNotationsSequenceValue()

	reqBody := `{
		"notation": "testnotation"
	}`

	DoSimpleTest(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{strconv.FormatUint(wordId, 10)},
		reqBody,
		nc.CreateNotation,
		http.StatusUnauthorized,
		"{}",
	)

	// DBにレコードが追加されない
	assert.Equal(t, 0, getCountFromNotations(notationId))
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

	wordId := insertIntoWords("りんご", "", 1)
	appleSentenceId := insertIntoSentences("赤い林檎を食べた", 1)
	lemonSentenceId := insertIntoSentences("黄色い檸檬を食べた", 1)

	reqBody := `{
		"notation": "林檎"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{strconv.FormatUint(wordId, 10)},
		reqBody,
		nc.CreateNotation,
	)

	// 「赤い林檎を食べた」には「林檎」が含まれるため、
	// sentences_wordsに追加される
	assert.Equal(t, 1, getCountFromSentencesWords(appleSentenceId, wordId))

	// 「黄色い檸檬を食べた」には「林檎」が含まれないため、
	// sentences_wordsに追加されない
	assert.Equal(t, 0, getCountFromSentencesWords(lemonSentenceId, wordId))
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

	wordId := insertIntoWords("りんご", "", 1)
	sentenceId := insertIntoSentences("赤い林檎を食べた", 2)

	reqBody := `{
		"notation": "林檎"
	}`

	ExecController(
		t,
		http.MethodPost,
		"/words/:wordId/notations",
		[]string{"wordId"},
		[]string{strconv.FormatUint(wordId, 10)},
		reqBody,
		nc.CreateNotation,
	)

	// 「赤い林檎を食べた」には「林檎」が含まれるが、
	// user_idが異なるため、
	// sentences_wordsに追加されない
	assert.Equal(t, 0, getCountFromSentencesWords(sentenceId, wordId))
}

func TestUpdateNotation(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを更新できることをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := insertIntoWords("test word", "test memo", 1)
	notationId := insertIntoNotations(wordId, "test notation")

	reqBody := `{
		"notation": "updated notation"
	}`

	// 更新されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %d,
			"word_id": %d,
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
		[]string{strconv.FormatUint(notationId, 10)},
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

	insertIntoWords("test word", "test memo", 1)

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

func TestUpdateNotation_UpdatedAssociation_AllNotationsDeleted(t *testing.T) {
	// Notationを更新した時、それによってSentence中にWordが含まれないことになったら、
	// sentences_wordsから値が削除されることをテスト

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

func TestUpdateNotation_UpdatedAssociation_SomeNotationRemaining(t *testing.T) {
	// Notationを更新した時、それでもSentence中にWordが含まれている場合、
	// sentences_wordsから値が削除されないことをテスト

	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := createTestWord(t, "りんご", "").Id
	createTestNotation(t, wordId, "林檎")
	notationId := createTestNotation(t, wordId, "リンゴ").Id
	sentenceId := createTestSentence(t, "林檎はリンゴと読むらしい").Id

	body := `{
		"notation": "RINGO"
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

	// Word「りんご」が、追加されていたNotationにより「林檎」「リンゴ」にもマッチしているので、
	// Notation「リンゴ」を「RINGO」に更新しても、Wordは「林檎」でマッチし続けるため、
	// sentences_wordsから削除されない
	assert.Equal(t, 1, getCountFromSentencesWords(sentenceId, wordId))
}

func TestDeleteNotation(t *testing.T) {
	// ログイン中のUserに紐づくWordに対し、Notationを削除できることをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := insertIntoWords("test word", "test memo", 1)
	notationId := insertIntoNotations(wordId, "test notation")

	// 削除されたレコードが返る
	expectedResponse := fmt.Sprintf(`
		{
			"id": %d,
			"word_id": %d,
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
		[]string{strconv.FormatUint(notationId, 10)},
		"",
		nc.DeleteNotation,
		http.StatusAccepted,
		expectedResponse,
	)

	// DBのレコードが削除される
	assert.Equal(t, 0, getCountFromNotations(notationId))
}

func TestDeleteNotation_WithInvalidUser(t *testing.T) {
	// ログイン中のUserに紐づかないWordに対し、Notationを削除できないことをテスト
	// TODO ログイン機能
	// とりあえずログインUserはuser_id=1とする
	DeleteAllFromWords()
	DeleteAllFromNotations()

	wordId := insertIntoWords("test word", "test memo", 2)
	notationId := insertIntoNotations(wordId, "test notation")
	assert.Equal(t, 1, getCountFromNotations(notationId))
	DoSimpleTest(
		t,
		http.MethodDelete,
		"/notations/:notationId",
		[]string{"notationId"},
		[]string{strconv.FormatUint(notationId, 10)},
		"",
		nc.DeleteNotation,
		http.StatusUnauthorized,
		"{}",
	)

	// DBのレコードが削除されない
	assert.Equal(t, 1, getCountFromNotations(notationId))
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
	// Notation「林檎」を削除しても、Wordは「リンゴ」でマッチし続けるため、
	// sentences_wordsから削除されない
	assert.Equal(t, 1, getCountFromSentencesWords(sentenceId, wordId))
}