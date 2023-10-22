package test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// 複数のテストで共通して使うパターンを切り出しまとめたファイル

func DoSimpleTest(
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
