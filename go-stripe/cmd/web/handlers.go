package main

import (
	"errors"
	"go-stripe/internal/cards"
	"go-stripe/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

//displays the virtual terminal page
func (app *application) virtualHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit the handler")

	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}

}

//displays the home page
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit the handler")

	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}

}

type TransactionData struct {
	FirstName       string
	LastName        string
	CardHolder      string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

//Get transaction data from post and stripe and returns the valid transactiondata or possibly an error
func (app *application) GetTransactionData(r *http.Request) (TransactionData, error) {

	var txnData TransactionData
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	//read posted data
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	cardHolder := r.Form.Get("cardholder_name")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	email := r.Form.Get("email")

	amount, err := strconv.Atoi(paymentAmount)

	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}
	card := cards.Card{
		Secret: app.Config.stripe.secret,
		Key:    app.Config.stripe.key,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err

	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err

	}
	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	txnData = TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		CardHolder:      cardHolder,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		PaymentAmount:   amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     int(expiryMonth),
		ExpiryYear:      int(expiryYear),
		BankReturnCode:  pi.Charges.Data[0].ID,
	}

	return txnData, nil
}

func (app *application) paymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	txnData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	widgetId, err := strconv.Atoi(r.Form.Get("product_id"))
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	//create a new customer
	customerId, err := app.SaveCustomer(txnData.FirstName, txnData.LastName, txnData.Email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Println(customerId)
	//create a new transaction

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		TransactionStatusID: 2,
	}
	txnId, err := app.SaveTransaction(txn)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Println(txnId)

	//create a new order
	var order models.Order
	order.WidgetID = widgetId
	order.TransactionID = txnId
	order.CustomerId = customerId
	order.StatusID = 1
	order.Quantity = 1
	order.Amount = txn.Amount
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err = app.SaveOrder(order)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// data := make(map[string]any)
	// data["cardholder"] = txnData.CardHolder
	// data["email"] = txnData.Email
	// data["pi"] = txnData.PaymentIntentID
	// data["pa"] = txnData.PaymentAmount
	// data["pc"] = txnData.PaymentCurrency
	// data["pm"] = txnData.PaymentMethodID
	// data["last_four"] = txnData.LastFour
	// data["expiry_month"] = txnData.ExpiryMonth
	// data["expiry_year"] = txnData.ExpiryYear
	// data["bank_return_code"] = txnData.BankReturnCode
	// data["first_name"] = txnData.FirstName
	// data["last_name"] = txnData.LastName

	//you should writethis data into session and redirect to another page to avoid recharging the credit card on reload
	// if err := app.renderTemplate(w, r, "succeeded", &templateData{Data: data}); err != nil {
	// 	app.errorLog.Println(err)
	// 	return
	// }

	//write transactiondata to session
	app.SessionManager.Put(r.Context(), "receipt", txnData)
	http.Redirect(w, r, "/receipt", http.StatusSeeOther)
}

func (app *application) Receipt(w http.ResponseWriter, r *http.Request) {
	txn := app.SessionManager.Get(r.Context(), "receipt").(TransactionData)
	data := make(map[string]any)
	data["txn"] = txn
	app.SessionManager.Remove(r.Context(), "receipt")
	if err := app.renderTemplate(w, r, "receipt", &templateData{Data: data}); err != nil {
		app.errorLog.Println(err)
		return
	}
}

//VirtualTerminalPaymentSuccedded for virtual terminal transactions
func (app *application) VirtualTerminalPaymentSuccedded(w http.ResponseWriter, r *http.Request) {

	txnData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	//create a new transaction

	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		TransactionStatusID: 2,
	}
	_, err = app.SaveTransaction(txn)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	//write transactiondata to session
	app.SessionManager.Put(r.Context(), "virtualTerminalReceipt", txnData)
	http.Redirect(w, r, "/virtual-terminal-payment-succeeded", http.StatusSeeOther)
}

//VirtualTerminalReceipt
func (app *application) VirtualTerminalReceipt(w http.ResponseWriter, r *http.Request) {
	txn := app.SessionManager.Get(r.Context(), "virtualTerminalReceipt").(TransactionData)
	data := make(map[string]any)
	data["txn"] = txn
	app.SessionManager.Remove(r.Context(), "virtualTerminalReceipt")
	if err := app.renderTemplate(w, r, "virtual-terminal-receipt", &templateData{Data: data}); err != nil {
		app.errorLog.Println(err)
		return
	}
}

//Gets info for the widget customer wants to pay for and renders the gohtml page
func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)
	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["widget"] = widget
	if err := app.renderTemplate(w, r, "buy-once", &templateData{Data: data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
		return
	}
}

//gets informations about the bronze subscription plan and renders the gohtml page
func (app *application) BronzePlan(w http.ResponseWriter, r *http.Request) {
	widget, err := app.DB.GetWidget(2)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := make(map[string]any)
	data["widget"] = widget
	// intMap := make(map[string]string)
	// intMap["plan_id"] = widget.PlanID
	if err := app.renderTemplate(w, r, "bronze-plan", &templateData{Data: data}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
		return
	}
}
func (app *application) BronzePlanReceipt(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "receipt-plan", &templateData{}); err != nil {
		app.errorLog.Println(err)
		return
	}
}

//Saves customer and returns Id
func (app *application) SaveCustomer(firstname, lastname, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil

}

//Saves transactions and returns the transaction Id
func (app *application) SaveTransaction(txn models.Transaction) (int, error) {

	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//Saves order and returns the orde Id
func (app *application) SaveOrder(ord models.Order) (int, error) {

	id, err := app.DB.InsertOrder(ord)
	if err != nil {
		return 0, err
	}
	return id, nil

}
func (app *application) ValidateFormData(data string) (string, error) {
	if data == "" {
		app.errorLog.Println("Form Data cannot be empty")
		return "", errors.New("invalid form input")
	}
	return data, nil
}
