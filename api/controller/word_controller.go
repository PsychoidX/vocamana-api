package controller

import (
	"api/usecase"
	"api/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type IWordController interface {
	GetAllWords(c echo.Context) error
	GetWordById(c echo.Context) error
	CreateWord(c echo.Context) error
	DeleteWord(c echo.Context) error
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
	return c.JSON(http.StatusCreated, wordResponses)
}

func (wc *WordController) GetWordById(c echo.Context) error {
	// id := c.Param("wordId")
	// TODO
	return nil
}

func (wc *WordController) CreateWord(c echo.Context) error {
	var wordInput model.WordRegistrationInput

	if err := c.Bind(&wordInput); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordRes, err := wc.wu.CreateWord(wordInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, wordRes)
}

func (wc *WordController) DeleteWord(c echo.Context) error {
	var req model.WordDeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordRes, err := wc.wu.DeleteWord(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusAccepted, wordRes)
}