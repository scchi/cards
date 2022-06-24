package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var testApp = &application{}

var jsonResponse struct {
	Id        int  `json:"deck_id"`
	Shuffled  bool `json:"shuffled"`
	Remaining int  `json:"remaining"`
}

func TestGetDeck(t *testing.T) {
	t.Run("returns 200 status code", func(t *testing.T) {
		handler := testApp.showDeckHandler
		router := httprouter.New()
		router.GET("/v1/decks/:id", handler)

		request, _ := http.NewRequest(http.MethodGet, "/v1/decks/5", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		json.NewDecoder(response.Body).Decode(&jsonResponse)

		gotId := jsonResponse.Id
		wantId := 5

		if gotId != wantId {
			t.Errorf("Wrong id")
		}

		if status := response.Code; status != http.StatusOK {
			t.Errorf("Wrong status")
		}
	})
}
