package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/scchi/cards/internal/assert"
	"github.com/scchi/cards/internal/data"
)

type createBody struct {
	Cards    []string `json:"cards"`
	Shuffled bool     `json:shuffled"`
}

var errorResponse struct {
	Error string `json:"error"`
}

var deck data.Deck

type cardsArray struct {
	Cards []struct {
		Value string `json:"value"`
		Suit  string `json:"suit"`
		Code  string `json:"code"`
	} `json:"cards"`
}

func TestCreateDeck(t *testing.T) {
	app := newTestApplication(t)
	path := "/v1/decks"

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("returns http.StatusUnprocessableEntity for invalid cards values", func(t *testing.T) {
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

			statusCode, _, _ := ts.post(t, path, bytes.NewReader(js))
			assert.Equal(t, statusCode, http.StatusUnprocessableEntity)
		}
	})

	t.Run("Returns http.StatusBadRequest for body with wrong value types", func(t *testing.T) {
		testBodies := []map[string]interface{}{
			{
				"shuffled": "stringShuffled",
			},
			{
				"shuffled": 1,
			},
			{
				"cards": true,
			},
			{
				"cards": "cardString",
			},
		}

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, _ := ts.post(t, path, bytes.NewReader(js))
			assert.Equal(t, statusCode, http.StatusBadRequest)
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
				"shuffled": true,
				"cards":    []string{"AC", "KH", "6D"},
			},
		}

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, _ := ts.post(t, path, bytes.NewReader(js))
			assert.Equal(t, statusCode, http.StatusCreated)
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

			_, header, _ := ts.post(t, path, bytes.NewReader(js))

			locationHeader := header["Location"][0]
			uuid := locationHeader[len(path)+1:]

			if !strings.HasPrefix(locationHeader, want) {
				t.Errorf("expecting header prefix of %s, got %s", want, locationHeader)
			}
			assert.Equal(t, len(uuid), 36)
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
				Shuffled: true,
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

			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

			gotRemaining := deck.Remaining
			var wantRemaining int
			if testBody.Cards == nil || len(testBody.Cards) == 0 {
				wantRemaining = 52
			} else {
				wantRemaining = len(testBody.Cards)
			}

			gotShuffled := deck.Shuffled
			wantShuffled := testBody.Shuffled

			gotID := deck.ID
			wantID := data.MockID

			assert.Equal(t, gotID, wantID)
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
		statusCode, _, _ := ts.get(t, fmt.Sprintf("/v1/decks/%s", data.MockID))
		assert.Equal(t, statusCode, http.StatusOK)
	})

	t.Run("Returns deck_id, remaining, and shuffled for valid id", func(t *testing.T) {
		_, _, body := ts.get(t, fmt.Sprintf("/v1/decks/%s", data.MockID))
		json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

		assert.Equal(t, deck.ID, data.MockID)
		assert.Equal(t, deck.Shuffled, data.MockShuffled)
		assert.Equal(t, deck.Remaining, len(data.MockCards))
	})

	t.Run("Returns cards array for valid id with each card having suit, value and code fields", func(t *testing.T) {
		_, _, body := ts.get(t, fmt.Sprintf("/v1/decks/%s", data.MockID))

		var ca cardsArray
		json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

		firstCard := ca.Cards[0]
		gotSuit := firstCard.Suit
		gotValue := firstCard.Value
		gotCode := firstCard.Code

		// TODO: hard-coded values
		assert.Equal(t, gotSuit, "SPADES")
		assert.Equal(t, gotValue, "ACE")
		assert.Equal(t, gotCode, "AS")
	})
}

func TestDrawDeck(t *testing.T) {
	app := newTestApplication(t)

	t.Run("Should return an error if count is less than 1 or greater than 52", func(t *testing.T) {
		counts := []map[string]int{
			{
				"count": 0,
			},
			{
				"count": -1,
			},
			{
				"count": 53,
			},
		}

		for _, count := range counts {
			countBytes, _ := json.Marshal(count)
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/v1/decks/%s", data.MockID), bytes.NewReader(countBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			app.routes().ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
		}
	})

	t.Run("Should return an error if deck has been drawn", func(t *testing.T) {
		countBytes, _ := json.Marshal(map[string]int{"count": 3})
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/v1/decks/%s", data.MockID), bytes.NewReader(countBytes))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		app.routes().ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
	})

	t.Run("Returns http.StatusNotFound for invalid id", func(t *testing.T) {
		countBytes, _ := json.Marshal(map[string]int{"count": 1})
		req, err := http.NewRequest(http.MethodPut, "/v1/decks/wrongid", bytes.NewReader(countBytes))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		app.routes().ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusNotFound)
	})

	t.Run("Should return an error if deck with given id doesn't exist", func(t *testing.T) {
		countBytes, _ := json.Marshal(map[string]int{"count": 1})
		req, err := http.NewRequest(http.MethodPut, "/v1/decks/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", bytes.NewReader(countBytes))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		app.routes().ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusNotFound)
	})

	t.Run("Should return JSON with cards array and with each card having suit, value, and code fields", func(t *testing.T) {
		counts := []int{
			1, 2,
		}

		for _, count := range counts {
			countBytes, _ := json.Marshal(map[string]int{"count": count})
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/v1/decks/%s", data.MockID), bytes.NewReader(countBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			app.routes().ServeHTTP(rr, req)

			var ca cardsArray
			json.NewDecoder(rr.Body).Decode(&ca)

			firstCard := ca.Cards[0]
			gotSuit := firstCard.Suit
			gotValue := firstCard.Value
			gotCode := firstCard.Code

			// TODO: hard-coded values
			assert.Equal(t, gotSuit, "SPADES")
			assert.Equal(t, gotValue, "ACE")
			assert.Equal(t, gotCode, "AS")

			if count == 2 {
				secondCard := ca.Cards[1]
				gotSuit = secondCard.Suit
				gotValue = secondCard.Value
				gotCode = secondCard.Code

				assert.Equal(t, gotSuit, "DIAMONDS")
				assert.Equal(t, gotValue, "9")
				assert.Equal(t, gotCode, "9D")
			}

			assert.Equal(t, len(ca.Cards), count)
			assert.Equal(t, rr.Code, http.StatusOK)
		}
	})

}
