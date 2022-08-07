package main

import (
	"net/http"
	// "myapp/models"
	"go-stripe/internal/models"
)

func (app *application) virtualHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit the handler")

	stringMap := make(map[string]string)

	stringMap["publishable_key"] = app.Config.stripe.key
	if err := app.renderTemplate(w, r, "terminal", &templateData{
		StringMap: stringMap,
	},"stripe-js"); err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) paymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	//read posted data
	cardHolder := r.Form.Get("cardholder_name")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	email := r.Form.Get("email")

	data := make(map[string]any)
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency
	data["pm"] = paymentMethod

	if err := app.renderTemplate(w, r, "succeeded", &templateData{Data: data}); err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {

	widget := models.Widget{
		ID:1,
		Name:"Custome Widget",
		Description:"This is a custom and very ni ce widget",
		InventoryLevel:10,
		Price:1000,
	}

	data := make(map[string]any)
	data["widget"] = widget
	if err := app.renderTemplate(w, r, "buy-once", &templateData{Data:data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
		return
	}
}
