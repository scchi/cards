package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/scchi/cards/internal/assert"
)

var testApp = &application{}

var jsonResponse struct {
	Id        int  `json:"deck_id"`
	Shuffled  bool `json:"shuffled"`
	Remaining int  `json:"remaining"`
}

// func TestGetDeck(t *testing.T) {
// 	t.Run("returns 200 status code", func(t *testing.T) {
// 		handler := testApp.showDeckHandler
// 		router := httprouter.New()
// 		router.GET("/v1/decks/:id", handler)

// 		request, _ := http.NewRequest(http.MethodGet, "/v1/decks/5", nil)
// 		response := httptest.NewRecorder()

// 		router.ServeHTTP(response, request)

// 		json.NewDecoder(response.Body).Decode(&jsonResponse)

// 		gotId := jsonResponse.Id
// 		wantId := 5

// 		if gotId != wantId {
// 			t.Errorf("Wrong id")
// 		}

// 		if status := response.Code; status != http.StatusOK {
// 			t.Errorf("Wrong status")
// 		}
// 	})
// }

type createBody struct {
	Cards    []string `json:"cards"`
	Shuffled bool     `json:shuffled"`
}

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

			gotId := jsonResponse.Id
			wantId := 0

			gotRemaining := jsonResponse.Remaining
			var wantRemaining int
			if testBody.Cards == nil || len(testBody.Cards) == 0 {
				wantRemaining = 52
			} else {
				wantRemaining = len(testBody.Cards)
			}

			gotShuffled := jsonResponse.Shuffled
			wantShuffled := testBody.Shuffled

			assert.Equal(t, gotId, wantId)
			assert.Equal(t, gotRemaining, wantRemaining)
			assert.Equal(t, gotShuffled, wantShuffled)
		}
	})
}

// var createResponse struct {
// 	Error map[string]string `json:"error"`
// }

// json.NewDecoder(response.Body).Decode(&createResponse)

// fmt.Println(createResponse.Error.cards)
