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

	words, err := wc.wu.GetAllWords(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var wordResponses []model.WordResponse
	for _, word := range words {
		wordRes := model.WordResponse{
			Id: word.Id,
			Word: word.Word,
			Memo: word.Memo,
			UserId: word.UserId,
		}
		wordResponses = append(wordResponses, wordRes)
	}

	return c.JSON(http.StatusOK, wordResponses)
}

func (wc *WordController) GetWordById(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	word, err := wc.wu.GetWordById(userId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(word == model.Word{}) {
		// usecaseで取得した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

	wordResponse := model.WordResponse{
		Id: word.Id,
		Word: word.Word,
		Memo: word.Memo,
		UserId: word.UserId,
	}
	return c.JSON(http.StatusOK, wordResponse)
}

func (wc *WordController) CreateWord(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	var req model.WordCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	WordCreation := model.WordCreation{
		Word: req.Word,
		Memo: req.Memo,
		UserId: userId,
	}

	word, err := wc.wu.CreateWord(WordCreation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	wordRes := model.WordResponse{
		Id:     word.Id,
		Word:   word.Word,
		Memo:   word.Memo,
		UserId: word.UserId,
	}

	return c.JSON(http.StatusCreated, wordRes)
}

func (wc *WordController) DeleteWord(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	word, err := wc.wu.DeleteWord(userId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(word == model.Word{}) {
		// usecaseでWordが削除されなかった場合
		// {}を返す
		return c.JSON(http.StatusAccepted, make(map[string]interface{}))
	}

	wordRes := model.WordResponse{
		Id:     word.Id,
		Word:   word.Word,
		Memo:   word.Memo,
		UserId: word.UserId,
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

	wordUpdate := model.WordUpdate{
		Id: wordId,
		Word: req.Word,
		Memo: req.Memo,
		UserId: userId,
	}

	word, err := wc.wu.UpdateWord(wordUpdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(word == model.Word{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusAccepted, make(map[string]interface{}))
	}

	wordRes := model.WordResponse{
		Id:     word.Id,
		Word:   word.Word,
		Memo:   word.Memo,
		UserId: word.UserId,
	}

	return c.JSON(http.StatusAccepted, wordRes)
}