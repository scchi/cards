package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/scchi/cards/internal/data"
	"github.com/scchi/cards/internal/validator"
)

func (app *application) createDeckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Shuffled bool     `json:"shuffled,omitempty"`
		Cards    []string `json:"cards"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	deck := &data.Deck{
		Shuffled:  input.Shuffled,
		Remaining: 5, // TODO
		Cards:     input.Cards,
	}

	v := validator.New()

	if data.ValidateDeck(v, deck); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Add Remaining field here

	err = app.models.Decks.Insert(deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", deck.ID))

	err = app.writeJSON(w, http.StatusCreated, deck, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
	}

	// TODO: Custom JSON Encoder for Cards

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
