package main

import (
	"flag"
	"fmt"
	"go-stripe/internal/driver"
	"go-stripe/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type Config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	Config   Config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func main() {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 4001, "server port to listen")
	flag.StringVar(&cfg.env, "environment", "development", "Application Environment {development || production || maintainancce}")
	flag.StringVar(&cfg.db.dsn, "dsn", "funbi:beedayme@tcp(localhost:3306)/widget?parseTime=true&tls=false", "DSN")

	flag.Parse()
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "API-INFO👹\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "API-ERROR👹\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDb(cfg.db.dsn)
	if err != nil {
		log.Println("Cannot connect to mysql 😞")
		errorLog.Fatal(err)
		return
	}

	defer conn.Close()
	app := &application{
		Config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) Serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.Config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting the backend  server in %s on port %d", app.Config.env, app.Config.port)
	return srv.ListenAndServe()
}
