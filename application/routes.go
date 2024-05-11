package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rexcfnghk/pricing-store/handler"
	"github.com/rexcfnghk/pricing-store/repository/currencymapping"
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/providercurrencyconfig"
	"github.com/rexcfnghk/pricing-store/repository/quote"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/providers", a.loadProviderRoutes)

	a.router = router
}

func (a *App) loadProviderRoutes(router chi.Router) {
	quoteRedisRepo := &quote.RedisRepo{
		Client: a.rdb,
	}

	providerRedisRepo := &provider.RedisRepo{
		Client: a.rdb,
	}
	currencyMappingRedisRepo := &currencymapping.RedisRepo{
		Client: a.rdb,
	}
	providerCurrencyConfigRedisRepo := &providercurrencyconfig.RedisRepo{
		Client: a.rdb,
	}

	quoteHandler := &handler.Quote{
		QuoteRepo:           quoteRedisRepo,
		ProviderRepo:        providerRedisRepo,
		CurrencyMappingRepo: currencyMappingRedisRepo,
	}

	providerHandler := &handler.Provider{
		ProviderRepo:               providerRedisRepo,
		CurrencyMappingRepo:        currencyMappingRedisRepo,
		ProviderCurrencyConfigRepo: providerCurrencyConfigRedisRepo,
	}

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(a.tokenAuth))
		r.Use(jwtauth.Authenticator(a.tokenAuth))
		r.Get("/bestprice", providerHandler.GetBestPrice)
	})

	router.Post("/{id}/quotes", quoteHandler.Create)

	router.Get("/{id}/providercurrencyconfigs", providerHandler.GetCurrencyConfigByCurrencyPair)
	router.Put("/{id}/providercurrencyconfigs", providerHandler.PutCurrencyConfigByCurrencyPair)
}
