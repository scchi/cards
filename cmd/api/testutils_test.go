package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/scchi/cards/internal/data"
	"github.com/scchi/cards/internal/jsonlog"
)

type createBody struct {
	Shuffled bool     `json:"shuffled"`
	Cards    []string `json:"cards"`
}

type cardsArray struct {
	Cards []struct {
		Value string `json:"value"`
		Suit  string `json:"suit"`
		Code  string `json:"code"`
	} `json:"cards"`
}

func newTestApplication(t *testing.T) *application {
	return &application{
		logger: jsonlog.New(io.Discard, jsonlog.LevelInfo),
		models: data.NewMockModels(),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) post(t *testing.T, urlPath string, testBody io.Reader) (int, http.Header, []byte) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "json", testBody)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, body
}
func put(t *testing.T, app *application, body map[string]int, url string) (*httptest.ResponseRecorder, *http.Request) {
	countBytes, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(countBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.routes().ServeHTTP(rr, req)
	return rr, req
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://test:passwOrd@localhost/test")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("../../migrations/test/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("../../migrations/test/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}

func countDecks(t *testing.T, db *sql.DB, count *int) {
	query := "SELECT COUNT(*) FROM decks"

	err := db.QueryRow(query).Scan(count)
	if err != nil {
		t.Fatal(err)
	}
}

func rowExists(t *testing.T, db *sql.DB, id string, count *int) {
	query := `
	SELECT COUNT(*)
	FROM decks
	WHERE id::text = $1`

	err := db.QueryRow(query, id).Scan(count)
	if err != nil {
		t.Fatal(err)
	}
}

func getSuit(card string) string {
	return data.Card(card).GetSuit()
}

func getValue(card string) string {
	return data.Card(card).GetValue()
}
