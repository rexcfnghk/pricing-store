package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rexcfnghk/pricing-store/handler"
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/quote"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/providers", a.loadQuoteRoutes)

	a.router = router
}

func (a *App) loadQuoteRoutes(router chi.Router) {
	quoteHandler := &handler.Quote{
		QuoteRepo: &quote.RedisRepo{
			Client: a.rdb,
		},
		ProviderRepo: &provider.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/{id}/quotes", quoteHandler.Create)
}
