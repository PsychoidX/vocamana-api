package test

import (
	"api/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func DeleteAllFromWords() {
	// wordsテーブルのレコードを全件削除
	db.Exec("TRUNCATE TABLE words CASCADE;")
	// word_id_seqシーケンスを1にリセット
	// nextval()で、2から連番で取得される
	db.Exec("SELECT setval('word_id_seq', 1);")
}

func GetCurrentWordsSequenceValue() uint64 {
	var currval uint64
	db.QueryRow(
		"SELECT currval('word_id_seq');",
	).Scan(&currval)
	return currval
}

func GetNextWordsSequenceValue() uint64 {
	// インデックスのカウンタを進めず参照のみするための実装
	return GetCurrentWordsSequenceValue() + 1
}

func DeleteAllFromSentences() {
	// sentencesテーブルのレコードを全件削除
	db.Exec("TRUNCATE TABLE sentences CASCADE;")
	db.Exec("SELECT setval('sentence_id_seq', 1);")
}

func GetCurrentSentencesSequenceValue() uint64 {
	var currval uint64
	db.QueryRow(
		"SELECT currval('sentence_id_seq');",
	).Scan(&currval)
	return currval
}

func GetNextSentencesSequenceValue() uint64 {
	// インデックスのカウンタを進めず参照のみするための実装
	return GetCurrentSentencesSequenceValue() + 1
}

func DeleteAllFromNotations() {
	// wordsテーブルのレコードを全件削除
	db.Exec("TRUNCATE TABLE notations CASCADE;")
	// word_id_seqシーケンスを1にリセット
	// nextval()で、2から連番で取得される
	db.Exec("SELECT setval('notation_id_seq', 1);")
}

func GetCurrentNotationsSequenceValue() uint64 {
	var currval uint64
	db.QueryRow(
		"SELECT currval('notation_id_seq');",
	).Scan(&currval)
	return currval
}

func GetNextNotationsSequenceValue() uint64 {
	// インデックスのカウンタを進めず参照のみするための実装
	return GetCurrentNotationsSequenceValue() + 1
}

func toMap(rec *httptest.ResponseRecorder) map[string]interface{} {
	// io.ReadAll(rec.Body)を使うと、内部でrec.Body.Readが呼ばれ、バッファが解放される
	// これにより、ReadAllでは最初の1回しかリクエストボディを取得できないため、Bytes()を使用
	var bodyMap map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &bodyMap)
	return bodyMap
}

func getBodyValueFromRecorder(rec *httptest.ResponseRecorder, key string) string {
	// recに記録されたリクエストボディ内の、keyの値を取得
	// interface型のbody[key]をstring型に変換
	return fmt.Sprintf("%v", toMap(rec)[key])
}

func toSentenceResponse(rec *httptest.ResponseRecorder) model.SentenceResponse {
	bodyMap := toMap(rec)
	id := fmt.Sprintf("%v", bodyMap["id"])
	sentence := fmt.Sprintf("%v", bodyMap["sentence"])
	userId := fmt.Sprintf("%v", bodyMap["user_id"])

	intId, _ := strconv.ParseUint(id, 10, 32)
	intUserId, _ := strconv.ParseUint(userId, 10, 32)

	return model.SentenceResponse{
		Id: intId,
		Sentence: sentence,
		UserId: intUserId,
	}
}

func toWordResponse(rec *httptest.ResponseRecorder) model.WordResponse {
	bodyMap := toMap(rec)
	id := fmt.Sprintf("%v", bodyMap["id"])
	word := fmt.Sprintf("%v", bodyMap["word"])
	memo := fmt.Sprintf("%v", bodyMap["memo"])
	userId := fmt.Sprintf("%v", bodyMap["user_id"])

	intId, _ := strconv.ParseUint(id, 10, 32)
	intUserId, _ := strconv.ParseUint(userId, 10, 32)

	return model.WordResponse{
		Id: intId,
		Word: word,
		Memo: memo,
		UserId: intUserId,
	}
}

func createTestSentence(t *testing.T, sentence string) model.SentenceResponse {
	// CreateSentenceを呼び出す
	// 他メソッドのテスト用データを作る用途で使用
	body := fmt.Sprintf(`
			{
				"sentence": "%s"
			}
		`,
		sentence,
	)

	_, rec := ExecController(
		t,
		"/sentences",
		sc.CreateSentence,
		HttpMethod(http.MethodPost),
		Body(body),
	)

	return toSentenceResponse(rec)
}

func createTestWord(t *testing.T, word, memo string) model.WordResponse {
	// CreateWordを呼び出す
	// 他メソッドのテスト用データを作る用途で使用
	body := fmt.Sprintf(`
			{
				"word": "%s",
				"memo": "%s"
			}
		`,
		word,
		memo,
	)

	_, rec := ExecController(
		t,
		"/words",
		wc.CreateWord,
		HttpMethod(http.MethodPost),
		Body(body),
	)

	return toWordResponse(rec)
}

func getCountFromSentencesWords[T uint64|string](sentenceId, wordId T) int {
	var count int

	db.QueryRow(`
		SELECT COUNT(*) FROM sentences_words
		WHERE sentence_id = $1
			AND word_id = $2;
		`,
		sentenceId,
		wordId,
	).Scan(&count)

	return count
}

func toNotationResponse(rec *httptest.ResponseRecorder) model.NotationResponse {
	bodyMap := toMap(rec)
	id := fmt.Sprintf("%v", bodyMap["id"])
	wordId := fmt.Sprintf("%v", bodyMap["word_id"])
	notation := fmt.Sprintf("%v", bodyMap["notation"])

	intId, _ := strconv.ParseUint(id, 10, 32)
	intWordId, _ := strconv.ParseUint(wordId, 10, 32)
	
	return model.NotationResponse{
		Id: intId,
		WordId: intWordId,
		Notation: notation,
	}
}

func createTestNotation(t *testing.T, wordId uint64, notation string) model.NotationResponse {
	// CreateNotationを呼び出す
	// 他メソッドのテスト用データを作る用途で使用
	body := fmt.Sprintf(`
			{
				"notation": "%s"
			}
		`,
		notation,
	)

	_, rec := ExecController(
		t,
		"/words/:wordId/notations",
		nc.CreateNotation,
		HttpMethod(http.MethodPost),
		Params(
			[]string{"wordId"},
			[]string{strconv.FormatUint(wordId, 10)},
		),
		Body(body),
	)

	return toNotationResponse(rec)	
}

func getCountFromNotations[T uint64|string](notationId T) int {
	var count int

	db.QueryRow(`
		SELECT COUNT(*) FROM notations
		WHERE id = $1
		`,
		notationId,
	).Scan(&count)

	return count
}

func getCountFromNotationsByNotation[T uint64|string](wordId T, notation string) int {
	var count int

	db.QueryRow(`
		SELECT COUNT(*) FROM notations
		WHERE word_id = $1
			AND notation = $2
		`,
		wordId,
		notation,
	).Scan(&count)

	return count
}


func insertIntoWords(word, memo string, userId uint64) uint64 {
	var wordId uint64
	db.QueryRow(`
		INSERT INTO words
		(id, word, memo, user_id)
		VALUES(nextval('word_id_seq'), $1, $2, $3)
		RETURNING id;
		`,
		word,
		memo,
		userId,
	).Scan(&wordId)

	return wordId
}

func insertIntoNotations(wordId uint64, notation string) uint64 {
	var notationId uint64
	db.QueryRow(`
		INSERT INTO notations
		(id, word_id, notation)
		VALUES(nextval('notation_id_seq'), $1, $2)
		RETURNING id;
		`,
		wordId,
		notation,
	).Scan(&notationId)

	return notationId
}

func insertIntoSentences(sentence string, userId uint64) uint64 {
	var sentenceId uint64
	db.QueryRow(`
		INSERT INTO sentences
		(id, sentence, user_id)
		VALUES(nextval('sentence_id_seq'), $1, $2)
		RETURNING id;
		`,
		sentence,
		userId,
	).Scan(&sentenceId)

	return sentenceId
}