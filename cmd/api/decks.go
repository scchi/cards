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

		if data.ValidateDeck(v, deck); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
	}

	prepForInsert(deck)

	err = app.models.Decks.Insert(deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	prepForCreateResponse(deck)

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

	prepForShowResponse(deck)

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

	if input.Count <= 0 || input.Count > 52 {
		error := map[string]string{
			"count": "must be between one and 52",
		}

		app.failedValidationResponse(w, r, error)
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

	if len(deck.StringCards) == 0 {
		app.deckErrorResponse(w, r)
		return
	}

	if input.Count > len(deck.StringCards) {
		error := errors.New("count must not be more than remaining cards in deck")

		app.badRequestResponse(w, r, error)
		return
	}

	returnCards := deck.StringCards[:input.Count]
	updateCards := deck.StringCards[input.Count:]

	deck.Cards = data.GenerateCards(updateCards)

	err = app.models.Decks.Update(deck)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	returnDeck := data.GenerateCards(returnCards)

	err = app.writeJSON(w, http.StatusOK, map[string][]data.Card{"cards": returnDeck}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
