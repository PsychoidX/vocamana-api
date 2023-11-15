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
	CreateMultipleWords(c echo.Context) error
	DeleteWord(c echo.Context) error
	UpdateWord(c echo.Context) error
	GetAssociatedSentences(c echo.Context) error
	GetAssociatedSentencesWithLink(c echo.Context) error
}

type WordController struct {
	wu *usecase.WordUsecase
}

func NewWordController(wu *usecase.WordUsecase) IWordController {
	return &WordController{wu}
}

func (wc *WordController) GetAllWords(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	words, err := wc.wu.GetAllWords(loginUserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var wordResponses []model.WordResponse
	for _, word := range words {
		wordRes := model.WordResponse{
			Id:     word.Id,
			Word:   word.Word,
			Memo:   word.Memo,
			UserId: word.UserId,
		}
		wordResponses = append(wordResponses, wordRes)
	}

	return c.JSON(http.StatusOK, wordResponses)
}

func (wc *WordController) GetWordById(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	word, err := wc.wu.GetWordById(loginUserId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if (word == model.Word{}) {
		// usecaseで取得した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

	wordRes := model.WordResponse{
		Id:     word.Id,
		Word:   word.Word,
		Memo:   word.Memo,
		UserId: word.UserId,
	}
	return c.JSON(http.StatusOK, wordRes)
}

func (wc *WordController) CreateWord(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.WordCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	WordCreation := model.WordCreation{
		Word:   req.Word,
		Memo:   req.Memo,
		LoginUserId: loginUserId,
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

func (wc *WordController) CreateMultipleWords(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	// TODO: words[]の中に不適切な形式のデータが入っていた場合、
	// すべてを登録失敗とするのではなく、不適切なデータのみを弾く実装にする
	
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.MultipleWordsCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var wordCreations []model.WordCreation
	for _, wordCreationReq := range req.Words {
		wordCreation := model.WordCreation{
			Word: wordCreationReq.Word,
			Memo: wordCreationReq.Memo,
			LoginUserId:   loginUserId,
		}
		wordCreations = append(wordCreations, wordCreation)
	} 

	words, err := wc.wu.CreateMultipleWords(wordCreations)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var wordResponses []model.WordResponse
	for _, word := range words {
		wordRes := model.WordResponse{
			Id:       word.Id,
			Word: word.Word,
			Memo: word.Memo,
			UserId: word.UserId,
		}
		wordResponses = append(wordResponses, wordRes)
	}

	return c.JSON(http.StatusCreated, wordResponses)
}

func (wc *WordController) DeleteWord(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	word, err := wc.wu.DeleteWord(loginUserId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if (word == model.Word{}) {
		// usecaseでWordが削除されなかった場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
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
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.WordUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordUpdate := model.WordUpdate{
		Id:     wordId,
		Word:   req.Word,
		Memo:   req.Memo,
		LoginUserId: loginUserId,
	}

	word, err := wc.wu.UpdateWord(wordUpdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if (word == model.Word{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	wordRes := model.WordResponse{
		Id:     word.Id,
		Word:   word.Word,
		Memo:   word.Memo,
		UserId: word.UserId,
	}

	return c.JSON(http.StatusAccepted, wordRes)
}

func (wc *WordController) GetAssociatedSentences(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentences, err := wc.wu.GetAssociatedSentencesByWordId(loginUserId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var sentenceResponses []model.SentenceResponse
	for _, sentence := range sentences {
		sentenceRes := model.SentenceResponse{
			Id:       sentence.Id,
			Sentence: sentence.Sentence,
			UserId:   sentence.UserId,
		}
		sentenceResponses = append(sentenceResponses, sentenceRes)
	}

	return c.JSON(http.StatusOK, sentenceResponses)
}

func (wc *WordController) GetAssociatedSentencesWithLink(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceWithLinks, err := wc.wu.GetAssociatedSentencesWithLinkByWordId(loginUserId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var sentenceWithLinkResponses []model.SentenceWIthLinkResponse
	for _, sentenceWithLink := range sentenceWithLinks {
		sentenceWithLinkRes := model.SentenceWIthLinkResponse{
			Id:               sentenceWithLink.Id,
			SentenceWithLink: sentenceWithLink.SentenceWithLink,
			UserId:           sentenceWithLink.UserId,
		}

		sentenceWithLinkResponses = append(sentenceWithLinkResponses, sentenceWithLinkRes)
	}

	return c.JSON(http.StatusOK, sentenceWithLinkResponses)
}
