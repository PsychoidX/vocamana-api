package main

import (
	"api/controller"
	"api/db"
	"api/repository"
	"api/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := db.NewDB()

	
	wr := repository.NewWordRepository(db)
	sr := repository.NewSentenceRepository(db)
	swr := repository.NewSentencesWordsRepository(db)
	nr := repository.NewNotationRepository(db)

	w := e.Group("/words")
	wu := usecase.NewWordUsecase(wr)
	wc := controller.NewWordController(wu)
	w.GET("", wc.GetAllWords)
	w.GET("/:wordId", wc.GetWordById)
	w.POST("", wc.CreateWord)
	w.PUT("/:wordId", wc.UpdateWord)
	w.DELETE("/:wordId", wc.DeleteWord)
	
	s := e.Group("/sentences")
	su := usecase.NewSentenceUsecase(sr, wr, swr)
	sc := controller.NewSentenceController(su)
	s.GET("", sc.GetAllSentences)
	s.GET("/:sentenceId", sc.GetSentenceById)
	s.POST("", sc.CreateSentence)
	s.PUT("/:sentenceId", sc.UpdateSentence)
	s.DELETE("/:sentenceId", sc.DeleteSentence)

	sa := e.Group("/sentences/association")
	sa.POST("/:sentenceId", sc.AssociateSentenceWithWords)

	n := e.Group("/words/:wordId/notations")
	nu := usecase.NewNotationUsecase(nr, wr)
	nc := controller.NewNotationController(nu)
	n.GET("", nc.GetAllNotations)
	n.POST("", nc.CreateNotation)
	n.PUT("/:notationId", nc.UpdateNotation)

	e.Logger.Fatal(e.Start(":8080"))
}