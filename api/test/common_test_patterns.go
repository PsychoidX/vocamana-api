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
	paramNames []string,
	paramValues []string,
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

	req := httptest.NewRequest(httpMethod, "/", bodyReader)

	// リクエストボディがある場合、JSON形式であるとする
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath(path)
	if paramNames != nil && paramValues != nil {
		c.SetParamNames(paramNames...)
		c.SetParamValues(paramValues...)
	}

	if assert.NoError(t, controllerMethod(c)) {
		assert.Equal(t, expectedStatusCode, rec.Code)
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}
