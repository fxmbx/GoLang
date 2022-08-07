package main

import (
	"flag"
	"fmt"
	"go-stripe/internal/driver"
	"go-stripe/internal/models"

	// "go-stripe/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const (
	version    = "1.0.0"
	cssVersion = "1"
	port       = "8080"
)

var session *scs.SessionManager

type Config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	Config        Config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func main() {

	var cfg Config
	flag.IntVar(&cfg.port, "port", 4000, "server port to listen")
	flag.StringVar(&cfg.env, "environment", "development", "Application Environment {development || production}")
	flag.StringVar(&cfg.api, "Api", "http://localhost:4001", "Api Url")

	flag.StringVar(&cfg.db.dsn, "dsn", "funbi:beedayme@tcp(localhost:3306)/widget?parseTime=true&tls=false", "DSN")

	// flag.StringVar(&cfg.db.dsn, "dsn", "funbi:beedayme@tcp(localhost:3306)/widget?parseTime=true&tls=false", "DSN")

	flag.Parse()
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDb(cfg.db.dsn)
	if err != nil {
		log.Println("Cannot connect to mysql 😞")
		errorLog.Fatal(err)
		return
	}
	defer conn.Close()

	//set up session manager
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	tc := make(map[string]*template.Template)

	app := &application{
		Config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
		Session:       session,
	}

	err = app.Serve()
	if err != nil {
		app.errorLog.Panicln(err)
		// log.Println(err)
		panic(err)
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

	app.infoLog.Printf("Starting http server in %s on port %d", app.Config.env, app.Config.port)
	return srv.ListenAndServe()
}
