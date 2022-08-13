package main

import (
	"encoding/json"
	"fmt"
	"go-stripe/internal/cards"
	"go-stripe/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v72"
)

type jsonResponse struct {
	OK      bool          `json:"ok"`
	Message string        `json:"message,omitempty"`
	Content string        `json:"content,omitempty"`
	ID      int           `json:"id,omitempty"`
	Token   *models.Token `json:"authentication_token,omitempty"`
}

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	CardBrand     string `json:"card_brand"`
	ExpiryYear    int    `json:"expiry_year"`
	ExpiryMonth   int    `json:"expiry_month"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	ProductID     string `json:"product_id"`
}

type Authenticate struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
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

func (app *application) GetWidgetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.Marshal(widget)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) CreateCustomerAndSubscribeToPlane(w http.ResponseWriter, r *http.Request) {
	var data stripePayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Println(data.Email, data.LastFour, data.Email, data.PaymentMethod)
	card := cards.Card{
		Secret:   app.Config.stripe.secret,
		Key:      app.Config.stripe.key,
		Currency: data.Currency,
	}

	okay := true
	var subscription *stripe.Subscription
	txnMsg := "Transaction Successful"

	stripeCustomer, msg, err := card.CreateCustomer(data.FirstName, data.LastName, data.PaymentMethod, data.Email)
	if err != nil {
		app.errorLog.Println(err)
		okay = false
		txnMsg = msg
	}

	if okay {
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.Email, data.LastFour, "")
		if err != nil {
			app.errorLog.Println(err)
			okay = false
			msg = "Error subscribing customer"
		}
		app.infoLog.Println("sub id is: ", subscription.ID)
	}

	if okay {
		product_Id, err := strconv.Atoi(data.ProductID)
		if err != nil {
			app.errorLog.Println("could not convert id to int : \n", err)
			return
		}

		customer_id, err := app.SaveCustomer(data.FirstName, data.LastName, data.Email)
		if err != nil {
			app.errorLog.Println("could not save customer : \n", err)
			return
		}

		//creat new transaction
		amount, _ := strconv.Atoi(data.Amount)

		txn := models.Transaction{
			Amount:              amount,
			ExpiryMonth:         data.ExpiryMonth,
			ExpiryYear:          data.ExpiryYear,
			Currency:            "ngn",
			LastFour:            data.LastFour,
			PaymentMethod:       data.PaymentMethod,
			TransactionStatusID: 2,
		}
		txnID, err := app.SaveTransaction(txn)
		if err != nil {
			app.errorLog.Println("could not save customer : \n", err)
			return
		}

		//create an order
		order := models.Order{
			WidgetID:      product_Id,
			TransactionID: txnID,
			CustomerId:    customer_id,
			Amount:        amount,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Quantity:      1,
			StatusID:      1,
		}
		_, err = app.SaveOrder(order)
		if err != nil {
			app.errorLog.Println("could not save customer : \n", err)
			return
		}
	}

	response := jsonResponse{
		OK:      okay,
		Message: txnMsg,
	}

	out, err := json.Marshal(response)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func (app *application) Authenticate(w http.ResponseWriter, r *http.Request) {

}

func (app *application) CreateAUthToken(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"passowrd"`
	}
	if err := app.readJSON(w, r, &userInput); err != nil {
		app.errorJson(w, err)
		return
	}
	app.infoLog.Println("Email from user input is ", userInput.Email)
	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w, "eml")
		// app.errorJson(w, err)
		return
	}
	validPassword, err := app.passwordMatched(userInput.Password, user.Password)
	if err != nil {
		app.invalidCredentials(w, "psd")
		return
		// app.errorLog.Println(err)
	}
	if !validPassword {
		app.invalidCredentials(w, "psd")
		return
		// app.errorLog.Println(err)
	}
	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.errorJson(w, err)
	}
	if err := app.DB.InsertToken(token, user); err != nil {
		app.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.OK = true
	payload.Message = "Success"
	payload.Content = fmt.Sprintf("token for %s generated successfully", userInput.Email)
	payload.Token = token

	app.writeJson(w, http.StatusAccepted, payload)

}

//Saves customer and returns
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

func (app *application) SaveTransaction(txn models.Transaction) (int, error) {

	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (app *application) SaveOrder(order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil

}
