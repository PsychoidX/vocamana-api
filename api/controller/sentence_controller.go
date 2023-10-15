package controller

import (
	"api/usecase"
	"api/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ISentenceController interface {
	GetAllSentences(c echo.Context) error
	CreateSentence(c echo.Context) error
}

type SentenceController struct {
	su *usecase.SentenceUsecase
}

func NewSentenceController(su *usecase.SentenceUsecase) ISentenceController {
	return &SentenceController{su}
}

func (sc *SentenceController) GetAllSentences(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得
	
	sentenceResponses, err := sc.su.GetAllSentences(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, sentenceResponses)
}

func (sc *SentenceController) CreateSentence(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	var req model.SentenceCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceRes, err := sc.su.CreateSentence(userId, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusCreated, sentenceRes)
}