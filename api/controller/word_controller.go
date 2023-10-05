package controller

import (
	"api/usecase"
	"github.com/labstack/echo/v4"
)

type IWordController interface {
	GetAllWords(c echo.Context) error
	GetWordById(c echo.Context) error
}

type WordController struct {
	wu *usecase.WordUsecase
}

func NewWordController(wu *usecase.WordUsecase) IWordController {
	return &WordController{wu}
}

func (wc *WordController) GetAllWords(c echo.Context) error {
	// TODO
	return nil
}

func (wc *WordController) GetWordById(c echo.Context) error {
	// id := c.Param("wordId")
	// TODO
	return nil
}
