package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/scchi/cards/internal/data"
)

func (app *application) createDeckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "create a new deck")
}

func (app *application) showDeckHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := app.readIDParam(ps)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	deck := data.Deck{
		ID:        id,
		Shuffled:  true,
		Remaining: 50,
		Cards: []data.Card{
			data.Card{
				Value: "ACE",
				Suit:  "DIAMOND",
				Code:  "AD",
			},
		},
		CreatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, deck, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) drawCardsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := app.readIDParam(ps)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	count, err := app.readCountParam(ps)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "draw %d cards from deck %d\n", count, id)
}
