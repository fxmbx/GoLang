package main

import (
	"encoding/json"
	"go-stripe/internal/cards"
	"net/http"
	"strconv"
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
	var payload stripePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	// err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorLog.Println(err)
		// app.errorJson(w, err)
		return
	}
	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.Config.stripe.secret,
		Key:      app.Config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true

	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		// err = app.writeJson(w, http.StatusAccepted, pi)
		out, err := json.Marshal(pi)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}
		// err = app.writeJson(w, http.StatusAccepted, j)
		out, err := json.Marshal(j)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)

	}

}
