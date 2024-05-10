package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rexcfnghk/pricing-store/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/quotes", loadQuoteRoutes)

	return router
}

func loadQuoteRoutes(router chi.Router) {
	quoteHandler := &handler.Quote{}

	router.Post("/", quoteHandler.Create)
}
