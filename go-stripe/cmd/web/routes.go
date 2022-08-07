package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	// mux.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"https://*", "http://*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))

	// mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/virtual-terminal", app.virtualHandler)
	mux.Post("/payment-succeeded", app.paymentSucceeded)

	mux.Get("/charge-once", app.ChargeOnce)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
