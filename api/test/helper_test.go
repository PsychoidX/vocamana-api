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

func GetCurrentWordsSequenceValue() int {
	var currval int
	db.QueryRow(
		"SELECT currval('word_id_seq');",
	).Scan(&currval)
	return currval
}

func GetNextWordsSequenceValue() int {
	// インデックスのカウンタを進めず参照のみするための実装
	return GetCurrentWordsSequenceValue() + 1
}

func DeleteAllFromSentences() {
	// sentencesテーブルのレコードを全件削除
	db.Exec("TRUNCATE TABLE sentences CASCADE;")
	db.Exec("SELECT setval('sentence_id_seq', 1);")
}

func GetCurrentSentencesSequenceValue() int {
	var currval int
	db.QueryRow(
		"SELECT currval('sentence_id_seq');",
	).Scan(&currval)
	return currval
}

func GetNextSentencesSequenceValue() int {
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

func GetCurrentNotationsSequenceValue() int {
	var currval int
	db.QueryRow(
		"SELECT currval('notation_id_seq');",
	).Scan(&currval)
	return currval
}

func GetNextNotationsSequenceValue() int {
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
		http.MethodPost,
		"/sentences",
		nil,
		nil,
		body,
		sc.CreateSentence,
	)

	return toSentenceResponse(rec)
}