package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/scchi/cards/internal/data"
	"github.com/scchi/cards/internal/validator"
)

func (app *application) createDeckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Shuffled bool        `json:"shuffled"`
		Cards    []data.Card `json:"cards"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	deck := &data.Deck{
		Shuffled: input.Shuffled,
		Cards:    input.Cards,
	}

	if deck.Cards != nil {
		v := validator.New()

		if data.ValidateCardsInput(v, deck); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
	}

	app.prepForInsert(deck)

	err = app.models.Decks.Insert(deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.prepForCreateResponse(deck)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/decks/%s", deck.ID))

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

	deck, err := app.models.Decks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.prepForShowResponse(deck)

	err = app.writeJSON(w, http.StatusOK, deck, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) drawCardsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Count int `json:"count"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if app.validateCount(v, input.Count); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	id, err := app.readIDParam(ps)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	deck, err := app.models.Decks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	v = validator.New()

	if app.validateForDraw(v, input.Count, deck); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cardsForReturn := deck.StringCards[:input.Count]
	cardsForUpdate := deck.StringCards[input.Count:]

	deck.Cards = data.GenerateCards(cardsForUpdate)
	returnCards := data.GenerateCards(cardsForReturn)

	err = app.models.Decks.Update(deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, map[string][]data.Card{"cards": returnCards}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
