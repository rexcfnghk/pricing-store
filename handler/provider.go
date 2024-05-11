package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/rexcfnghk/pricing-store/repository/currencypair"
	"github.com/rexcfnghk/pricing-store/repository/customer"
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/providercurrencyconfig"
	"github.com/rexcfnghk/pricing-store/repository/quote"
)

type Provider struct {
	ProviderRepo               *provider.RedisRepo
	CurrencyPairRepo           *currencypair.RedisRepo
	ProviderCurrencyConfigRepo *providercurrencyconfig.RedisRepo
	CustomerRepo               *customer.RedisRepo
	QuoteRepo                  *quote.RedisRepo
}

func (h *Provider) GetCurrencyConfigByCurrencyPair(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	marketProviderId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.ProviderRepo.GetById(r.Context(), marketProviderId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	base, quote := r.URL.Query().Get("base"), r.URL.Query().Get("quote")
	if base == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if quote == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currencyPairId, err := h.CurrencyPairRepo.GetByCurrencyPairId(r.Context(), base, quote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	providercurrencyconfig, err := h.ProviderCurrencyConfigRepo.GetById(r.Context(), marketProviderId, currencyPairId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(providercurrencyconfig)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

func (h *Provider) PutCurrencyConfigByCurrencyPair(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	marketProviderId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	base, quote := r.URL.Query().Get("base"), r.URL.Query().Get("quote")
	if base == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if quote == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.ProviderRepo.GetById(r.Context(), marketProviderId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var body model.ProviderCurrencyConfig
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currencyPairId, err := h.CurrencyPairRepo.GetByCurrencyPairId(r.Context(), base, quote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProviderCurrencyConfigRepo.UpdateById(r.Context(), marketProviderId, currencyPairId, body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Provider) GetBestPrice(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	customerId, err := strconv.Atoi(fmt.Sprintf("%s", claims["sub"]))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	base, quote := r.URL.Query().Get("base"), r.URL.Query().Get("quote")
	if base == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if quote == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currencyPairId, err := h.CurrencyPairRepo.GetByCurrencyPairId(r.Context(), base, quote)
	// TODO: Handle no quotes
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	quotes, err := h.QuoteRepo.GetAllByCurrencyPairId(r.Context(), currencyPairId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("%v", quotes)

	//customer, err := h.CustomerRepo.GetById(r.Context(), customerId)

	// Get customer rating factor from customer ID
	// Get currency mapping ID from query["base"] and query["quote"]
	// Get all quotes with "quotes:{currencymappingid}"
	// Get all currency configs with quotes.DistinctBy(q => q.MarketProviderId)
	// Filter quotes to only show active based on currency configs
	// BEST PRICE = max bid price and min ask price
}
