package main

import (
	"api/controller"
	"api/db"
	"api/repository"
	"api/usecase"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	db := db.NewDB()
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{
				os.Getenv("FE_URL"),
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			AllowHeaders: []string{},
		},
	))

	// Repository
	wr := repository.NewWordRepository(db)
	sr := repository.NewSentenceRepository(db)
	swr := repository.NewSentencesWordsRepository(db)
	nr := repository.NewNotationRepository(db)

	// Usecase
	wu := usecase.NewWordUsecase(wr, sr, swr, nr)
	su := usecase.NewSentenceUsecase(sr, wr, swr, nr)
	nu := usecase.NewNotationUsecase(wr, sr, swr, nr)

	// Controller
	wc := controller.NewWordController(wu)
	sc := controller.NewSentenceController(su)
	nc := controller.NewNotationController(nu)

	w := e.Group("/words")
	w.GET("", wc.GetAllWords)
	w.GET("/:wordId", wc.GetWordById)
	w.POST("", wc.CreateWord)
	w.POST("/multiple", wc.CreateMultipleWords)
	w.PUT("/:wordId", wc.UpdateWord)
	w.DELETE("/:wordId", wc.DeleteWord)
	w.GET("/:wordId/associated-sentences", wc.GetAssociatedSentencesWithLink)

	s := e.Group("/sentences")
	s.GET("", sc.GetAllSentences)
	s.GET("/:sentenceId", sc.GetSentenceById)
	s.POST("", sc.CreateSentence)
	s.POST("/multiple", sc.CreateMultipleSentences)
	s.PUT("/:sentenceId", sc.UpdateSentence)
	s.DELETE("/:sentenceId", sc.DeleteSentence)
	s.GET("/:sentenceId/associated-words", sc.GetAssociatedWords)

	sa := e.Group("/sentences/association")
	sa.POST("/:sentenceId", sc.AssociateSentenceWithWords)

	wn := e.Group("/words/:wordId/notations")
	wn.GET("", nc.GetAllNotations)
	wn.POST("", nc.CreateNotation)

	n := e.Group("/notations")
	n.PUT("/:notationId", nc.UpdateNotation)
	n.DELETE("/:notationId", nc.DeleteNotation)

	e.Logger.Fatal(e.Start(":8080"))
}
