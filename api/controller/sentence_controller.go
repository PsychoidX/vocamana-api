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
	
	sentences, err := sc.su.GetAllSentences(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var sentenceResponses []model.SentenceResponse
	for _, sentence := range sentences {
		sentenceRes := model.SentenceResponse{
			Id: sentence.Id,
			Sentence: sentence.Sentence,
			UserId: sentence.UserId,
		}
		sentenceResponses = append(sentenceResponses, sentenceRes)
	}

	return c.JSON(http.StatusOK, sentenceResponses)
}

func (sc *SentenceController) GetSentenceById(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentence, err := sc.su.GetSentenceById(userId, sentenceId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(sentence == model.Sentence{}) {
		// usecaseで取得した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}
	
	sentenceRes := model.SentenceResponse{
		Id: sentence.Id,
		Sentence: sentence.Sentence,
		UserId: sentence.UserId,
	}
	return c.JSON(http.StatusOK, sentenceRes)
}

func (sc *SentenceController) CreateSentence(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	var req model.SentenceCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceCreation := model.SentenceCreation{
		Sentence: req.Sentence,
		UserId: userId,
	}

	sentence, err := sc.su.CreateSentence(sentenceCreation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceRes := model.SentenceResponse{
		Id: sentence.Id,
		Sentence: sentence.Sentence,
		UserId: sentence.UserId,
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

	sentenceUpdate := model.SentenceUpdate{
		Id: sentenceId,
		Sentence: req.Sentence,
		UserId: userId,
	}

	sentence, err := sc.su.UpdateSentence(sentenceUpdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	if(sentence == model.Sentence{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}
	
	sentenceRes := model.SentenceResponse{
		Id: sentence.Id,
		Sentence: sentence.Sentence,
		UserId: sentence.UserId,
	}

	return c.JSON(http.StatusAccepted, sentenceRes)
}

func (sc *SentenceController) DeleteSentence(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentence, err := sc.su.DeleteSentence(userId, sentenceId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(sentence == model.Sentence{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	sentenceRes := model.SentenceResponse{
		Id: sentence.Id,
		Sentence: sentence.Sentence,
		UserId: sentence.UserId,
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

	wordIds := model.WordIds{
		WordIds: req.WordIds,
	}

	resultWordIds, err := sc.su.AssociateSentenceWithWords(userId, sentenceId, wordIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordIdsRes := model.WordIdsResponse{
		WordIds: resultWordIds.WordIds,
	}

	return c.JSON(http.StatusAccepted, wordIdsRes)
}
