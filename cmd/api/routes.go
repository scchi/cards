package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.GET("/v1/healthcheck", app.healthcheckHandler)
	router.POST("/v1/decks", app.createDeckHandler)
	router.GET("/v1/decks/:id", app.showDeckHandler)
	router.PUT("/v1/decks/:id", app.drawCardsHandler)
	return router
}
