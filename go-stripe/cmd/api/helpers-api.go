package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// type jsonResponse struct {
// 	Error   bool   `json:"error"`
// 	Message string `json:"message"`
// 	Data    any    `json:"data,omitempty"`
// }

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //1mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return app.errorJson(w, err)
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return app.errorJson(w, errors.New("body must have only a single JSON value"))
	}
	return nil
}

func (app *application) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		// app.errorLog.Println(err)
		return err
	}
	return nil

}

func (app *application) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.OK = false
	payload.Content = err.Error()
	payload.Message = "Error payload 😞 "
	log.Printf("\n\nError payload 😞 %v\n\n", payload)
	app.errorLog.Println(payload)

	return app.writeJson(w, statusCode, payload)
}
