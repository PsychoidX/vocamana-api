package controller

import (
	"api/model"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type INotationController interface {
	GetAllNotations(c echo.Context) error
	CreateNotation(c echo.Context) error
	UpdateNotation(c echo.Context) error
	DeleteNotation(c echo.Context) error
}

type NotationController struct {
	wu *usecase.WordUsecase
}

func NewNotationController(wu *usecase.WordUsecase) INotationController {
	return &NotationController{wu}
}

func (nc *NotationController) GetAllNotations(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	notations, err := nc.wu.GetAllNotations(loginUserId, wordId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var notationResponses []model.NotationResponse
	for _, notation := range notations {
		notationRes := model.NotationResponse{
			Id: notation.Id,
			WordId: notation.WordId,
			Notation: notation.Notation,
		}
		notationResponses = append(notationResponses, notationRes)
	}

	return c.JSON(http.StatusOK, notationResponses)
}

func (nc *NotationController) CreateNotation(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.NotationCreationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	wordId, err := strconv.ParseUint(c.Param("wordId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	notationCreation := model.NotationCreation{
		WordId: wordId,
		Notation: req.Notation,
		LoginUserId: loginUserId,
	}

	notation, err := nc.wu.CreateNotation(notationCreation)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(notation == model.Notation{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	notationRes := model.NotationResponse{
		Id: notation.Id,
		WordId: notation.WordId,
		Notation: notation.Notation,
	}
	
	return c.JSON(http.StatusCreated, notationRes)
}

func (nc *NotationController) UpdateNotation(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req model.NotationUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	notationId, err := strconv.ParseUint(c.Param("notationId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	notationUpdate := model.NotationUpdate{
		Id: notationId,
		Notation: req.Notation,
		LoginUserId: loginUserId,
	}

	notation, err := nc.wu.UpdateNotation(notationUpdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(notation == model.Notation{}) {
		// usecaseで更新した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	notationRes := model.NotationResponse{
		Id: notation.Id,
		WordId: notation.WordId,
		Notation: notation.Notation,
	}
	
	return c.JSON(http.StatusAccepted, notationRes)
}

func (nc *NotationController) DeleteNotation(c echo.Context) error {
	loginUserId, err := GetLoginUserId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	notationId, err := strconv.ParseUint(c.Param("notationId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	notation, err := nc.wu.DeleteNotation(loginUserId, notationId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if(notation == model.Notation{}) {
		// usecaseで削除した結果がゼロ値の場合
		// {}を返す
		return c.JSON(http.StatusUnauthorized, make(map[string]interface{}))
	}

	notationRes := model.NotationResponse{
		Id: notation.Id,
		WordId: notation.WordId,
		Notation: notation.Notation,
	}
	
	return c.JSON(http.StatusAccepted, notationRes)
}