package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	j := jsonResponse{
		OK: true,
	}
	out, err := json.Marshal(j)
	if err != nil {
		// log.Println()
		app.errorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
