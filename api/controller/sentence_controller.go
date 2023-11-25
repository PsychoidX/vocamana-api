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
	CreateMultipleSentences(c echo.Context) error
	UpdateSentence(c echo.Context) error
	DeleteSentence(c echo.Context) error
	GetAssociatedWords(c echo.Context) error
}

type SentenceController struct {
	su *usecase.SentenceUsecase
	au *usecase.AssociationUsecase
}

func NewSentenceController(
	su *usecase.SentenceUsecase,
	au *usecase.AssociationUsecase,
) ISentenceController {
	return &SentenceController{su, au}
}

func (sc *SentenceController) GetAllSentences(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentencesWithLink, err := sc.au.GetAllSentencesWithLink(loginUserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var sentenceWithLinkResponses []model.SentenceWIthLinkResponse
	for _, sentenceWithLink := range sentencesWithLink {
		res := model.SentenceWIthLinkResponse{
			Id:       sentenceWithLink.Id,
			Sentence: sentenceWithLink.Sentence,
			SentenceWithLink: sentenceWithLink.SentenceWithLink,
			UserId:   sentenceWithLink.UserId,
		}
		sentenceWithLinkResponses = append(sentenceWithLinkResponses, res)
	}

	return c.JSON(http.StatusOK, sentenceWithLinkResponses)
}

func (sc *SentenceController) GetSentenceById(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentence, err := sc.su.GetSentenceById(loginUserId, sentenceId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if (sentence == model.Sentence{}) {
		// usecaseで取得した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

	sentenceRes := model.SentenceResponse{
		Id:       sentence.Id,
		Sentence: sentence.Sentence,
		UserId:   sentence.UserId,
	}
	return c.JSON(http.StatusOK, sentenceRes)
}

func (sc *SentenceController) CreateSentence(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.SentenceCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceCreation := model.SentenceCreation{
		Sentence: req.Sentence,
		LoginUserId:   loginUserId,
	}

	sentence, err := sc.su.CreateSentence(sentenceCreation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceRes := model.SentenceResponse{
		Id:       sentence.Id,
		Sentence: sentence.Sentence,
		UserId:   sentence.UserId,
	}

	return c.JSON(http.StatusCreated, sentenceRes)
}

func (sc *SentenceController) CreateMultipleSentences(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// TODO: words[]の中に不適切な形式のデータが入っていた場合、
	// すべてを登録失敗とするのではなく、不適切なデータのみを弾く実装にする
	var req model.MultipleSentencesCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var sentenceCreations []model.SentenceCreation
	for _, sentenceCreationReq := range req.Sentences {
		sentenceCreation := model.SentenceCreation{
			Sentence: sentenceCreationReq.Sentence,
			LoginUserId:   loginUserId,
		}
		sentenceCreations = append(sentenceCreations, sentenceCreation)
	} 

	sentences, err := sc.su.CreateMultipleSentences(sentenceCreations)
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

	return c.JSON(http.StatusCreated, sentenceResponses)
}

func (sc *SentenceController) UpdateSentence(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.SentenceUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceUpdate := model.SentenceUpdate{
		Id:       sentenceId,
		Sentence: req.Sentence,
		LoginUserId:   loginUserId,
	}

	sentence, err := sc.su.UpdateSentence(sentenceUpdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if (sentence == model.Sentence{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	// クエリパラメータ ?with-link=true の場合、
	// レスポンスをリンク付きSentenceにする
	if(c.QueryParam("with-link") == "true") {
		sentenceWithLink, err := sc.au.GetSentenceWithLinkById(loginUserId, sentenceId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		sentenceWithLinkRes := model.SentenceWIthLinkResponse{
			Id: sentenceWithLink.Id,
			Sentence: sentenceWithLink.Sentence,
			SentenceWithLink: sentenceWithLink.SentenceWithLink,
			UserId: sentenceWithLink.UserId,
		}

		return c.JSON(http.StatusAccepted, sentenceWithLinkRes)
	} else {
		sentenceRes := model.SentenceResponse{
			Id:       sentence.Id,
			Sentence: sentence.Sentence,
			UserId:   sentence.UserId,
		}

		return c.JSON(http.StatusAccepted, sentenceRes)
	}
}

func (sc *SentenceController) DeleteSentence(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentence, err := sc.su.DeleteSentence(loginUserId, sentenceId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if (sentence == model.Sentence{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	sentenceRes := model.SentenceResponse{
		Id:       sentence.Id,
		Sentence: sentence.Sentence,
		UserId:   sentence.UserId,
	}

	return c.JSON(http.StatusAccepted, sentenceRes)
}

func (sc *SentenceController) GetAssociatedWords(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sentenceId, err := strconv.ParseUint(c.Param("sentenceId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	words, err := sc.su.GetAssociatedWordsBySentenceId(loginUserId, sentenceId)
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
