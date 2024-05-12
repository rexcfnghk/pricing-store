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
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/providercurrencyconfig"
	"github.com/rexcfnghk/pricing-store/service"
)

type Provider struct {
	ProviderRepo               *provider.RedisRepo
	CurrencyPairRepo           *currencypair.RedisRepo
	ProviderCurrencyConfigRepo *providercurrencyconfig.RedisRepo
	BestPriceService           *service.BestPriceService
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

	providerCurrencyConfig, err := h.ProviderCurrencyConfigRepo.GetById(r.Context(), marketProviderId, currencyPairId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(providerCurrencyConfig)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		fmt.Println("failed to write response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	currencyPair := &model.CurrencyPair{
		Base:  base,
		Quote: quote,
	}

	_, err = h.BestPriceService.GetBestPrice(r.Context(), currencyPair, customerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Add response
}
