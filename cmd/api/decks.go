package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createDeckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new deck")
}

func (app *application) showDeckHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of deck %d\n", id)
}

func (app *application) drawCardsHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	count, err := strconv.ParseInt(params.ByName("count"), 10, 64)
	if err != nil || count < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "draw %d cards from deck %d\n", count, id)
}
