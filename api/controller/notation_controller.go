package controller

import (
	"api/model"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type INotationController interface {
	CreateNotation(c echo.Context) error
}

type NotationController struct {
	nu *usecase.NotationUsecase
}

func NewNotationController(nu *usecase.NotationUsecase) INotationController {
	return &NotationController{nu}
}

func (nc *NotationController) CreateNotation(c echo.Context) error {
	var userId uint64 = 1 // TODO セッションから取得

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
	}

	notation, err := nc.nu.CreateNotation(userId, notationCreation)
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