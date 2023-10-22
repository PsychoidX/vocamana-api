package test

import (
	"api/controller"
	"api/db"
	"api/repository"
	"api/usecase"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func doSimpleTest(
		t *testing.T,
		httpMethod string,
		path string,
		body string,
		controllerMethod func(echo.Context) error,
		expectedStatusCode int,
		expectedJSON string,
	) {
	e := echo.New()
	
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	
	req := httptest.NewRequest(httpMethod, path, bodyReader)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllerMethod(c)) {
		assert.Equal(t, expectedStatusCode, rec.Code)
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}

func TestGetAllWords(t *testing.T) {
	db := db.NewDB()
	wr := repository.NewWordRepository(db)
	wu := usecase.NewWordUsecase(wr)
	wc := controller.NewWordController(wu)
	
	// レコードが1つも無い場合、[]ではなくnullが返る
	doSimpleTest(
		t,
		http.MethodGet,
		"/words",
		"a",
		wc.GetAllWords,
		http.StatusOK,
		`null`,
	)
}