package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scchi/cards/internal/assert"
	"github.com/scchi/cards/internal/data"
)

type createRequest struct {
	Shuffled bool     `json:"shuffled"`
	Cards    []string `json:"cards"`
}

var ca cardsArray

var path = "/v1/decks"

func TestCreate(t *testing.T) {
	app := newTestApplication(t)

	testDB := newTestDB(t)
	app.models = data.NewModels(testDB)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Sending a valid request via POST /v1/decks/:id (where cards is present and not empty) creates a record in db and returns the appropriate response", func(t *testing.T) {
		testBodies := []createRequest{
			{
				Shuffled: true,
				Cards:    []string{"AC", "KH", "6D"},
			},
			{
				Cards: []string{"AC", "AD", "AH"},
			},
			{
				Shuffled: false,
				Cards:    []string{"AC", "KH", "6D", "9H"},
			},
		}

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, body := ts.post(t, path, bytes.NewReader(js))

			var deck data.Deck
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

			var rowCount int
			rowExists(t, testDB, deck.ID, &rowCount)

			assert.Equal(t, rowCount, 1)
			assert.Equal(t, statusCode, http.StatusCreated)
			assert.Equal(t, deck.Shuffled, testBody.Shuffled)
			assert.Equal(t, deck.Remaining, len(testBody.Cards))
		}
	})

	t.Run("Sending a valid request via POST /v1/decks/:id (where cards is not present or is empty) creates a record in db and returns the appropriate response", func(t *testing.T) {
		testBodies := []createRequest{
			{
				Shuffled: true,
			},
			{
				Cards: []string{},
			},
		}

		for _, testBody := range testBodies {
			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, body := ts.post(t, path, bytes.NewReader(js))

			var deck data.Deck
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

			var rowCount int
			rowExists(t, testDB, deck.ID, &rowCount)

			assert.Equal(t, rowCount, 1)
			assert.Equal(t, statusCode, http.StatusCreated)
			assert.Equal(t, deck.Shuffled, testBody.Shuffled)
			assert.Equal(t, deck.Remaining, 52)
		}
	})

	t.Run("Cards must be shuffled when shuffled is set to true", func(t *testing.T) {
		testBody := createRequest{
			Shuffled: true,
			Cards:    []string{"AC", "KH", "8C", "9D", "2C", "AH", "QS", "10D"},
		}

		js, err := json.Marshal(testBody)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 30; i++ {
			_, header, body := ts.post(t, path, bytes.NewReader(js))

			var deck data.Deck
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

			locationHeader := header["Location"][0]
			_, _, body = ts.get(t, locationHeader)
			json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

			firstCard := ca.Cards[0]

			if firstCard.Code != "AC" {
				break
			}

			if i == 29 {
				t.Error("Cards must be shuffled")
			}
		}
	})

	t.Run("Cards must not be shuffled when shuffled is set to false", func(t *testing.T) {
		testBody := createRequest{
			Shuffled: false,
			Cards:    []string{"AC", "KH", "8C", "9D", "2C", "AH", "QS", "10D"},
		}

		js, err := json.Marshal(testBody)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 30; i++ {
			_, header, body := ts.post(t, path, bytes.NewReader(js))

			var deck data.Deck
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

			locationHeader := header["Location"][0]
			_, _, body = ts.get(t, locationHeader)

			json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

			firstCard := ca.Cards[0]

			if firstCard.Code == "AH" {
				t.Error("Cards must not be shuffled")
			}
		}
	})

	t.Run("Cards must not be shuffled when shuffled is not in the request body", func(t *testing.T) {
		testBody := createRequest{
			Cards: []string{"AC", "KH"},
		}

		js, err := json.Marshal(testBody)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 20; i++ {
			_, header, body := ts.post(t, path, bytes.NewReader(js))

			var deck data.Deck
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

			locationHeader := header["Location"][0]
			_, _, body = ts.get(t, locationHeader)
			json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

			firstCard := ca.Cards[0]

			if firstCard.Code == "KH" {
				t.Error("Cards must not be shuffled")
			}
		}
	})

	t.Run("Number of DB rows should increase by one after each valid request", func(t *testing.T) {
		testBodies := []createRequest{
			{},
			{
				Shuffled: true,
			},
			{
				Cards: []string{},
			},
			{
				Cards: []string{"AS"},
			},
			{
				Cards: []string{"7D", "AH"},
			},
			{
				Shuffled: true,
				Cards:    []string{"AC", "KH", "6D"},
			},
		}

		var beforeCount int
		var afterCount int

		for _, testBody := range testBodies {
			countDecks(t, testDB, &beforeCount)

			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}
			_, _, _ = ts.post(t, path, bytes.NewReader(js))

			countDecks(t, testDB, &afterCount)

			assert.Equal(t, afterCount, beforeCount+1)
		}
	})

	t.Run("Number of DB rows should not increase after an invalid request", func(t *testing.T) {
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

		var beforeCount int
		var afterCount int

		for _, testBody := range testBodies {
			countDecks(t, testDB, &beforeCount)

			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}
			_, _, _ = ts.post(t, path, bytes.NewReader(js))

			countDecks(t, testDB, &afterCount)

			assert.Equal(t, afterCount, beforeCount)
		}
	})
}

func TestGet(t *testing.T) {
	app := newTestApplication(t)

	testDB := newTestDB(t)
	app.models = data.NewModels(testDB)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Returns ErrRecordNotFound for non-existent row and invalid id's", func(t *testing.T) {
		ids := []string{
			data.MockID,
			"invalidID",
			"5",
		}

		for _, id := range ids {
			statusCode, _, _ := ts.get(t, path+id)

			assert.Equal(t, statusCode, http.StatusNotFound)
		}

	})

	t.Run("Returns 200 status and appropriate response for rows that exist", func(t *testing.T) {
		testBody := createRequest{
			Shuffled: true,
			Cards:    []string{"AC", "KH"},
		}

		js, err := json.Marshal(testBody)
		if err != nil {
			t.Fatal(err)
		}

		_, header, body := ts.post(t, path, bytes.NewReader(js))

		var deck data.Deck
		json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

		var rowCount int
		rowExists(t, testDB, deck.ID, &rowCount)

		assert.Equal(t, rowCount, 1)

		locationHeader := header["Location"][0]
		statusCode, _, body := ts.get(t, locationHeader)

		json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

		assert.Equal(t, statusCode, http.StatusOK)
		assert.Equal(t, deck.Shuffled, testBody.Shuffled)
		assert.Equal(t, deck.Remaining, len(testBody.Cards))
	})

	t.Run("Returns cards with value, suit, and code fields", func(t *testing.T) {
		testBody := createRequest{
			Shuffled: false,
			Cards:    []string{"AC", "KH"},
		}

		js, err := json.Marshal(testBody)
		if err != nil {
			t.Fatal(err)
		}

		_, header, body := ts.post(t, path, bytes.NewReader(js))

		var deck data.Deck
		json.NewDecoder(bytes.NewReader(body)).Decode(&deck)

		locationHeader := header["Location"][0]
		_, _, body = ts.get(t, locationHeader)
		json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

		firstCard := ca.Cards[0]
		assert.Equal(t, firstCard.Code, "AC")
		assert.Equal(t, firstCard.Suit, "CLUBS")
		assert.Equal(t, firstCard.Value, "ACE")

		secondCard := ca.Cards[1]
		assert.Equal(t, secondCard.Code, "KH")
		assert.Equal(t, secondCard.Suit, "HEARTS")
		assert.Equal(t, secondCard.Value, "KING")
	})
}

func TestDraw(t *testing.T) {
	app := newTestApplication(t)

	testDB := newTestDB(t)
	app.models = data.NewModels(testDB)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

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
		req, err := http.NewRequest(http.MethodPut, "/v1/decks/"+data.MockID, bytes.NewReader(countBytes))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		app.routes().ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusNotFound)
	})

	t.Run("Should return Cards array with length equal to count value in request body", func(t *testing.T) {
		newDeck := createRequest{
			Shuffled: false,
			Cards:    []string{"AC", "KH", "QD", "3H", "5S", "AH", "KS"},
		}

		js, err := json.Marshal(newDeck)
		if err != nil {
			t.Fatal(err)
		}

		_, header, _ := ts.post(t, path, bytes.NewReader(js))
		locationHeader := header["Location"][0]

		counts := []int{
			3, 1,
		}

		for _, count := range counts {
			countBytes, _ := json.Marshal(map[string]int{"count": count})
			req, err := http.NewRequest(http.MethodPut, locationHeader, bytes.NewReader(countBytes))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			app.routes().ServeHTTP(rr, req)

			json.NewDecoder(rr.Body).Decode(&ca)
			assert.Equal(t, len(ca.Cards), count)

			for idx, card := range ca.Cards {
				wantCard := newDeck.Cards[idx]

				assert.Equal(t, card.Suit, getSuit(wantCard))
				assert.Equal(t, card.Value, getValue(wantCard))
				assert.Equal(t, card.Code, wantCard)
			}

			newDeck.Cards = newDeck.Cards[count:]
			assert.Equal(t, rr.Code, http.StatusOK)
		}
	})

	t.Run("Deck length should be decremented by count value in request", func(t *testing.T) {
		newDeck := createRequest{
			Shuffled: false,
			Cards:    []string{"AC", "KH", "QD", "3H", "5S"},
		}

		js, err := json.Marshal(newDeck)
		if err != nil {
			t.Fatal(err)
		}

		counts := []int{
			1, 3,
		}

		for _, count := range counts {
			_, header, body := ts.post(t, path, bytes.NewReader(js))
			locationHeader := header["Location"][0]
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)
			beforeRemaining := deck.Remaining

			countBytes, _ := json.Marshal(map[string]int{"count": count})
			req, err := http.NewRequest(http.MethodPut, locationHeader, bytes.NewReader(countBytes))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			app.routes().ServeHTTP(rr, req)

			_, _, body = ts.get(t, locationHeader)
			json.NewDecoder(bytes.NewReader(body)).Decode(&deck)
			afterRemaining := deck.Remaining
			assert.Equal(t, beforeRemaining-count, afterRemaining)
		}
	})

	t.Run("Cards in deck should be drawn starting from the zeroth index", func(t *testing.T) {
		newDeck := createRequest{
			Shuffled: false,
			Cards:    []string{"AC", "KH", "QD", "3H", "5S"},
		}

		js, err := json.Marshal(newDeck)
		if err != nil {
			t.Fatal(err)
		}

		_, header, _ := ts.post(t, path, bytes.NewReader(js))
		locationHeader := header["Location"][0]

		counts := []int{
			1, 3,
		}

		for _, count := range counts {

			countBytes, _ := json.Marshal(map[string]int{"count": count})
			req, err := http.NewRequest(http.MethodPut, locationHeader, bytes.NewReader(countBytes))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			app.routes().ServeHTTP(rr, req)

			_, _, body := ts.get(t, locationHeader)
			json.NewDecoder(bytes.NewReader(body)).Decode(&ca)

			assert.Equal(t, ca.Cards[0].Code, newDeck.Cards[count])

			newDeck.Cards = newDeck.Cards[count:]
		}
	})
}
