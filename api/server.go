package main

import (
	"github.com/labstack/echo/v4"
	"api/db"
	"api/repository"
	"api/usecase"
	"api/controller"
)

func main() {
	e := echo.New()
	db := db.NewDB()
	wr := repository.NewWordRepository(db)
	wu := usecase.NewWordUsecase(wr)
	wc := controller.NewWordController(wu)
	e.GET("/words", wc.GetAllWords)
	e.Logger.Fatal(e.Start(":8080"))
}