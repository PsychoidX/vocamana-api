package controller

import (
	"api/model"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IWordController interface {
	GetAllWords(c echo.Context) error
	GetWordById(c echo.Context) error
	CreateWord(c echo.Context) error
	DeleteWord(c echo.Context) error
	UpdateWord(c echo.Context) error
}

type WordController struct {
	wu *usecase.WordUsecase
}

func NewWordController(wu *usecase.WordUsecase) IWordController {
	return &WordController{wu}
}

func (wc *WordController) GetAllWords(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	wordResponses, err := wc.wu.GetAllWords(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, wordResponses)
}

func (wc *WordController) GetWordById(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordResponse, err := wc.wu.GetWordById(userId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(wordResponse == model.WordResponse{}) {
		// usecaseで取得した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	} else {
		return c.JSON(http.StatusOK, wordResponse)
	}
}

func (wc *WordController) CreateWord(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	var req model.WordCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordRes, err := wc.wu.CreateWord(userId, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, wordRes)
}

func (wc *WordController) DeleteWord(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordRes, err := wc.wu.DeleteWord(userId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusAccepted, wordRes)
}

func (wc *WordController) UpdateWord(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	var req model.WordUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	req.Id = wordId

	wordRes, err := wc.wu.UpdateWord(userId, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusAccepted, wordRes)
}