package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/scchi/cards/internal/assert"
	"github.com/scchi/cards/internal/data"
)

func TestInsert(t *testing.T) {
	path := "/v1/decks"
	app := newTestApplication(t)

	testDB := newTestDB(t)
	app.models = data.NewModels(testDB)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("DB rows should increase by one after each valid request", func(t *testing.T) {
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

		var beforeCount int
		var afterCount int
		query := "SELECT COUNT(*) FROM decks"

		for _, testBody := range testBodies {
			err := testDB.QueryRow(query).Scan(&beforeCount)
			if err != nil {
				t.Fatal(err)
			}

			js, err := json.Marshal(testBody)
			if err != nil {
				t.Fatal(err)
			}

			_, _, _ = ts.post(t, path, bytes.NewReader(js))

			err = testDB.QueryRow(query).Scan(&afterCount)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(beforeCount, afterCount)
			assert.Equal(t, afterCount, beforeCount+1)
			assert.Equal(t, 201, http.StatusCreated)
		}
	})
}
