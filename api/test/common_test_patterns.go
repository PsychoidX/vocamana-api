package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// 複数のテストで共通して使うパターンを切り出しまとめたファイル

// controllerの呼び出し方法を指定するオプション
type CallControllerOption struct {
	httpMethod string
	paramNames []string
	paramValues []string
	queryParamNames []string
	queryParamValues [][]string
	body string
}

// CallControllerOptionを破壊的に変更するメソッド
type CallControllerOptionBuildFunc func(*CallControllerOption)

func Params(paramNames, paramValues []string) CallControllerOptionBuildFunc {
	// https://hoge.com/:foo/:piyo のURLパターンに、https://hoge.com/abc/5 でアクセスしたい場合、
	// paramNames == ["foo", "piyo"]
	// paramValues == ["abc", "5"]
	return func(opt *CallControllerOption) {
		if(len(paramNames) == len(paramValues)) {
			opt.paramNames = paramNames
			opt.paramValues = paramValues
		}
	}
}

func QueryParams(queryParamNames []string, queryParamValues [][]string) CallControllerOptionBuildFunc {
	// https://hoge.comのURLパターンに、https://hoge.com?foo=abc&piyo=1&piyo=5 でアクセスしたい場合、
	// queryParamNames == ["foo", "piyo"]
	// queryParamValues == [["abc"], ["1", "5"]]
	return func(opt *CallControllerOption) {
		if(len(queryParamNames) == len(queryParamValues)) {
			opt.queryParamNames = queryParamNames
			opt.queryParamValues = queryParamValues
		}
	}
}

func HttpMethod(method string) CallControllerOptionBuildFunc {
	return func(opt *CallControllerOption) {
		opt.httpMethod = method
	}
}

func Body(body string) CallControllerOptionBuildFunc {
	return func(opt *CallControllerOption) {
		opt.body = body
	}
}


func DoSimpleTest(
	t *testing.T,
	path string,
	controllerMethod func(echo.Context) error,
	expectedStatusCode int,
	expectedJSON string,
	buildFuncs ...CallControllerOptionBuildFunc,
) {
	isNoError, rec := ExecController(
		t,
		path,
		controllerMethod,
		buildFuncs...
	)

	if isNoError {
		assert.Equal(t, expectedStatusCode, rec.Code)
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}

func ExecController(
	t *testing.T,
	path string,
	controllerMethod func(echo.Context) error,
	buildFuncs ...CallControllerOptionBuildFunc,
) (
	bool,
	*httptest.ResponseRecorder,
) {
	// 返り値の検証をせず、Controllerの呼び出しのみを実行
	e := echo.New()

	option := CallControllerOption{
		httpMethod: http.MethodGet, // HTTPメソッドが指定されていなければGETを使用
	}

	// 引数で指定されたオプションをoptionに反映
	for _, f := range buildFuncs{
		f(&option)
	}

	var bodyReader io.Reader
	if option.body != "" {
		bodyReader = strings.NewReader(option.body)
	}

	var req *http.Request
	if option.queryParamNames != nil && option.queryParamValues != nil {
		queryParams := make(url.Values)
		
		for i:=0; i<len(option.queryParamNames); i++ {
			queryParams[option.queryParamNames[i]] = option.queryParamValues[i]
		}
		
		req = httptest.NewRequest(option.httpMethod, "/?"+ queryParams.Encode(), bodyReader)
	} else {
		req = httptest.NewRequest(option.httpMethod, "/", bodyReader)
	}

	// リクエストボディがある場合、JSON形式であるとする
	if option.body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath(path)
	if option.paramNames != nil && option.paramValues != nil {
		c.SetParamNames(option.paramNames...)
		c.SetParamValues(option.paramValues...)
	}
	
	return assert.NoError(t, controllerMethod(c)), rec
}
