package main

import (
	"bytes"
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

func TestCreateDeck(t *testing.T) {
	handler := testApp.createDeckHandler
	router := httprouter.New()
	router.POST("/v1/decks", handler)

	t.Run("returns http.StatusUnprocessableEntity for malformed JSON body", func(t *testing.T) {
		testBodies := []map[string]interface{}{
			{
				"shuffled": true,
			},
			{
				"shuffled": false,
				"cards":    []string{"AS", "AS"},
			},
			{
				"cards": []string{"RR"},
			},
		}

		for _, body := range testBodies {
			body, _ := json.Marshal(body)
			request, _ := http.NewRequest(http.MethodPost, "/v1/decks", bytes.NewReader(body))
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if status := response.Code; status != http.StatusUnprocessableEntity {
				t.Errorf("Wrong status")
			}
		}
	})

	// var createResponse struct {
	// 	Error map[string]string `json:"error"`
	// }

	// json.NewDecoder(response.Body).Decode(&createResponse)

	// fmt.Println(createResponse.Error.cards)
}

// func main() {
// 	mcPostBody := map[string]interface{}{
// 		"question_text": "Is this a test post for MutliQuestion?",
// 	}
// 	body, _ := json.Marshal(mcPostBody)
// 	req, err := http.NewRequest("POST", "/questions/", bytes.NewReader(body))
// 	var m map[string]interface{}
// 	err = json.NewDecoder(req.Body).Decode(&m)
// 	req.Body.Close()
// 	fmt.Println(err, m)
// }
