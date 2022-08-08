package models

import (
	"context"
	"database/sql"
	"time"
)

//type for database connection values
type DBModel struct {
	DB *sql.DB
}

//wrapper for all models
type Models struct {
	DB DBModel
	// Widget Widget
}

//Returns a model type with databse connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

//Widget is the type for all widget
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// order is the type for all orders
type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerId    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

//status is the type for all order statuses
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

//transactionstatus is a type for all transaction status
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

//Transaction is the type for all transactions
type Transaction struct {
	ID                  int       `json:"id"`
	Currency            string    `json:"currency"`
	Amount              int       `json:"amount"`
	TransactionStatusID int       `json:"transaction_status_id"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"bank_return_code"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

//Users is the type for all users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

//Custome is the type for all customers
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

//returns a widget gotten by its ID or an error
func (m *DBModel) GetWidget(id int) (Widget, error) {
	//use context and set a reasonable timeout everytime you connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, `
		select
			id, name, description, inventory_level, price, coalesce(image, ''), 
			created_at, updated_at 
		from 
			widgets
		where id = ?`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}
	return widget, nil
}

//Inserts a new transacton and returns newly inserted transaction id or possibly an error
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into transactions 
			(amount, currency, last_four,bank_return_code, expiry_month, expiry_year, payment_intent, payment_method,
			transaction_status_id,created_at, updated_at) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?,?,?)
	`

	result, err := m.DB.ExecContext(
		ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		txn.TransactionStatusID, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}

//Inserts a new order and returns newly inserted order id or possibly an error
func (m *DBModel) InsertOrder(ord Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into orders 
			(widget_id, transaction_id, status_id,quantity,
				amount, customer_id, created_at, updated_at) 
		values (?, ?, ?, ?, ?, ?, ? ,?)
	`

	result, err := m.DB.ExecContext(
		ctx, stmt,
		ord.WidgetID, ord.TransactionID, ord.StatusID,
		ord.Quantity, ord.Amount, ord.CustomerId, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}

//Inserts a new customer and returns newly inserted customer id or possibly an error
func (m *DBModel) InsertCustomer(cust Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into customers  
			(first_name, last_name, email,created_at, updated_at) 
		values (?, ?, ?, ?, ?)
	`

	result, err := m.DB.ExecContext(
		ctx, stmt,
		cust.FirstName, cust.LastName, cust.Email,
		time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil

}
