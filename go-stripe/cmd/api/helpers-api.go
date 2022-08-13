package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// type jsonResponse struct {
// 	Error   bool   `json:"error"`
// 	Message string `json:"message"`
// 	Data    any    `json:"data,omitempty"`
// }

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		app.errorJson(w, errors.New("body must have only a single JSON value"))
		return err
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
	payload.Message = err.Error()
	payload.Content = "Error payload 😞 "
	log.Printf("\n\nError payload 😞 %v\n\n", payload)
	app.errorLog.Println(payload)

	return app.writeJson(w, statusCode, payload)
}

func (app *application) invalidCredentials(w http.ResponseWriter, msg string) error {
	var payload jsonResponse
	payload.OK = false
	payload.Message = "Invalid credentials"
	payload.Content = fmt.Sprintf("invalid Credentials paylaod %v😞", msg)

	if err := app.writeJson(w, http.StatusUnauthorized, payload); err != nil {
		return err
	}
	return nil
}

func (app *application) passwordMatched(password, hashPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	log.Println(err.Error())
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
	// app.errorLog.Println()
}

func (app *application) generateToken() {

}
