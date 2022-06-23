package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(ps httprouter.Params) (int64, error) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) readCountParam(ps httprouter.Params) (int64, error) {
	count, err := strconv.ParseInt(ps.ByName("count"), 10, 64)
	if err != nil || count < 1 {
		return 0, errors.New("invalid count parameter")
	}

	return count, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
