package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
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

func getBodyValueFromRecorder(rec *httptest.ResponseRecorder, key string) string {
	// recに記録されたリクエストボディ内の、keyの値を取得
	resp := rec.Result()
	body, _ := io.ReadAll(resp.Body)

	var bodyMap map[string]interface{}
	json.Unmarshal(body, &bodyMap)

	// interface型のbody[key]をstring型に変換
	return fmt.Sprintf("%v", bodyMap[key])
}