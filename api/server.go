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

	w := e.Group("/words")
	wr := repository.NewWordRepository(db)
	wu := usecase.NewWordUsecase(wr)
	wc := controller.NewWordController(wu)
	w.GET("", wc.GetAllWords)
	w.POST("", wc.CreateWord)
	w.PUT("/:wordId", wc.UpdateWord)
	w.DELETE("/:wordId", wc.DeleteWord)
	
	e.Logger.Fatal(e.Start(":8080"))
}