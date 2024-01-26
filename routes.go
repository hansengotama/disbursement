package main

import (
	disbursementhandler "github.com/hansengotama/disbursement/internal/handler/disbursement"
	"github.com/hansengotama/disbursement/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func initRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.MiddlewareAuth)
	r.Route("/disbursements", func(r chi.Router) {
		r.Post("/request", disbursementhandler.HandleRequestDisbursement)
	})

	return r
}
