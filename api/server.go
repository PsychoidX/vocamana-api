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
	w.GET("/:wordId", wc.GetWordById)
	w.POST("", wc.CreateWord)
	w.PUT("/:wordId", wc.UpdateWord)
	w.DELETE("/:wordId", wc.DeleteWord)
	
	s := e.Group("/sentences")
	sr := repository.NewSentenceRepository(db)
	su := usecase.NewSentenceUsecase(sr)
	sc := controller.NewSentenceController(su)
	s.GET("", sc.GetAllSentences)
	s.POST("", sc.CreateSentence)

	e.Logger.Fatal(e.Start(":8080"))
}