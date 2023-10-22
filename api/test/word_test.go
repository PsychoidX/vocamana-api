package test

import (
	"net/http"
	"testing"
)

func TestGetAllWords(t *testing.T) {
	DeleteAllFromWords()
	
	// レコードが1つも無い場合、[]ではなくnullが返る
	DoSimpleTest(
		t,
		http.MethodGet,
		"/words",
		"a",
		wc.GetAllWords,
		http.StatusOK,
		`null`,
	)
}