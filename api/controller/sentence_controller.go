package controller

import (
	"api/model"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ISentenceController interface {
	GetAllSentences(c echo.Context) error
	GetSentenceById(c echo.Context) error
	CreateSentence(c echo.Context) error
	UpdateSentence(c echo.Context) error
	DeleteSentence(c echo.Context) error
	AssociateSentenceWithWords(c echo.Context) error
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

	return c.JSON(http.StatusOK, sentenceResponses)
}

func (sc *SentenceController) GetSentenceById(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceResponse, err := sc.su.GetSentenceById(userId, sentenceId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(sentenceResponse == model.SentenceResponse{}) {
		// usecaseで取得した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	} else {
		return c.JSON(http.StatusOK, sentenceResponse)
	}
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

func (sc *SentenceController) UpdateSentence(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	var req model.SentenceUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	req.Id = sentenceId

	sentenceRes, err := sc.su.UpdateSentence(userId, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusAccepted, sentenceRes)
}

func (sc *SentenceController) DeleteSentence(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceRes, err := sc.su.DeleteSentence(userId, sentenceId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	return c.JSON(http.StatusAccepted, sentenceRes)
}

func (sc *SentenceController) AssociateSentenceWithWords(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.WordIdsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordIdsRes, err := sc.su.AssociateSentenceWithWords(userId, sentenceId, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusAccepted, wordIdsRes)
}
