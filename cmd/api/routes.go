package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/v1/healthcheck", app.healthcheckHandler)
	router.POST("/v1/decks", app.createDeckHandler) // TODO: shuffled and partial as payload
	router.GET("/v1/decks/:id", app.showDeckHandler)
	router.PATCH("/v1/decks/:id/count/:count", app.drawCardsHandler) // TODO: count as payload

	return router
}
