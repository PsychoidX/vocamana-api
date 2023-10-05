package controller

import (
	"api/usecase"
	"github.com/labstack/echo/v4"
)

type IWordController interface {
	GetAllWords(c echo.Context) error
	GetWordById(c echo.Context) error
}

type wordController struct {
	wu usecase.IWordUsecase
}

func NewWordController(wu usecase.IWordUsecase) IWordController {
	return &wordController{wu}
}

func (wc *wordController) GetAllWords(c echo.Context) error {
	// TODO
	return nil
}

func (wc *wordController) GetWordById(c echo.Context) error {
	// id := c.Param("wordId")
	// TODO
	return nil
}
