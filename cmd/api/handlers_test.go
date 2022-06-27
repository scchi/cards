package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/scchi/cards/internal/assert"
	"github.com/scchi/cards/internal/data"
)

var jsonResponse struct {
	Id        int  `json:"deck_id"`
	Shuffled  bool `json:"shuffled"`
	Remaining int  `json:"remaining"`
}

type createBody struct {
	Cards    []string `json:"cards"`
	Shuffled bool     `json:shuffled"`
}

var errorResponse struct {
	Error string `json:"error"`
}

var deck data.Deck

func TestCreate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("returns http.StatusUnprocessableEntity", func(t *testing.T) {
		testBodies := []map[string]interface{}{
			{
				"shuffled": false,
				"cards":    []string{"AS", "AS"},
			},
			{
				"cards": []string{"RR"},
			},
			{
				"cards": []string{"8D", "ZZ"},
			},
		}

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, _ := ts.post(t, "/v1/decks", bytes.NewReader(js))
			assert.Equal(t, statusCode, http.StatusUnprocessableEntity)
		}
	})

	t.Run("returns http.StatusCreated for valid request body", func(t *testing.T) {
		testBodies := []map[string]interface{}{
			{},
			{
				"shuffled": true,
			},
			{
				"cards": []string{},
			},
			{
				"cards": []string{"AS"},
			},
			{
				"cards": []string{"7D", "AH"},
			},
			{
				"shuffled": false,
				"cards":    []string{"AC", "KH"},
			},
		}

		want := http.StatusCreated

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, _ := ts.post(t, "/v1/decks", bytes.NewReader(js))
			assert.Equal(t, statusCode, want)
		}
	})

	t.Run("Returns valid Location header for valid requests", func(t *testing.T) {
		testBodies := []map[string]interface{}{
			{},
			{
				"shuffled": true,
			},
			{
				"cards": []string{},
			},
			{
				"cards": []string{"AS"},
			},
			{
				"cards": []string{"7D", "AH"},
			},
			{
				"shuffled": false,
				"cards":    []string{"AC", "KH"},
			},
		}

		want := "/v1/decks/"

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			_, header, _ := ts.post(t, "/v1/decks", bytes.NewReader(js))

			locationHeader := header["Location"][0]

			// TODO: test suffix for UUID regex
			if !strings.HasPrefix(locationHeader, want) {
				t.Errorf("expecting header prefix of %s, got %s", want, locationHeader)
			}
		}
	})

	t.Run("Valid requests return JSON with deck_id, remaining, and shuffled fields ", func(t *testing.T) {
		testBodies := []createBody{
			{},
			{
				Shuffled: true,
				Cards:    []string{},
			},
			{
				Shuffled: false,
				Cards: []string{
					"AS", "4D",
				},
			},
			{
				Cards: []string{
					"KC", "3H", "QS",
				},
			},
		}

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			_, _, body := ts.post(t, "/v1/decks", bytes.NewReader(js))

			json.NewDecoder(bytes.NewReader(body)).Decode(&jsonResponse)

			gotRemaining := jsonResponse.Remaining
			var wantRemaining int
			if testBody.Cards == nil || len(testBody.Cards) == 0 {
				wantRemaining = 52
			} else {
				wantRemaining = len(testBody.Cards)
			}

			gotShuffled := jsonResponse.Shuffled
			wantShuffled := testBody.Shuffled

			// assert.Equal(t, gotId, wantId)
			assert.Equal(t, gotRemaining, wantRemaining)
			assert.Equal(t, gotShuffled, wantShuffled)
		}
	})
}

func TestGetDeck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Returns http.StatusNotFound for invalid id", func(t *testing.T) {
		statusCode, _, _ := ts.get(t, "/v1/decks/wrongid")
		assert.Equal(t, statusCode, http.StatusNotFound)
	})

	t.Run("Returns error message in JSON body", func(t *testing.T) {
		_, _, body := ts.get(t, "/v1/decks/wrongid")
		json.NewDecoder(bytes.NewReader(body)).Decode(&errorResponse)

		want := "the requested resource could not be found"
		assert.Equal(t, errorResponse.Error, want)
	})

	t.Run("Returns http.StatusOK for valid id", func(t *testing.T) {
		statusCode, _, _ := ts.get(t, "/v1/decks/whatanid")
		assert.Equal(t, statusCode, http.StatusOK)
	})

	t.Run("Returns deck_id, remaining, and shuffled for valid id", func(t *testing.T) {
		_, _, body := ts.get(t, "/v1/decks/existingid")
		json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

		assert.Equal(t, deck.ID, "existingid")
		assert.Equal(t, deck.Shuffled, false)
		assert.Equal(t, deck.Remaining, 2)
	})

	t.Run("Returns cards array for valid id with each card having suit, value and code fields", func(t *testing.T) {
		_, _, body := ts.get(t, "/v1/decks/existingid")

		type cardsArray struct {
			Cards []struct {
				Value string `json:"value"`
				Suit  string `json:"suit"`
				Code  string `json:"code"`
			} `json:"cards"`
		}

		var ca cardsArray
		json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

		firstCard := ca.Cards[0]
		suit := firstCard.Suit
		value := firstCard.Value
		code := firstCard.Code

		assert.Equal(t, suit, "SPADES")
		assert.Equal(t, value, "ACE")
		assert.Equal(t, code, "AS")
	})
}
