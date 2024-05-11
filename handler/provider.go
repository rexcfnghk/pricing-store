package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/rexcfnghk/pricing-store/repository/currencymapping"
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/providercurrencyconfig"
)

type Provider struct {
	ProviderRepo               *provider.RedisRepo
	CurrencyMappingRepo        *currencymapping.RedisRepo
	ProviderCurrencyConfigRepo *providercurrencyconfig.RedisRepo
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

	currencyMappingId, err := h.CurrencyMappingRepo.GetByCurrencyPairId(r.Context(), base, quote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	providercurrencyconfig, err := h.ProviderCurrencyConfigRepo.GetById(r.Context(), marketProviderId, currencyMappingId)
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

	currencyMappingId, err := h.CurrencyMappingRepo.GetByCurrencyPairId(r.Context(), base, quote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProviderCurrencyConfigRepo.UpdateById(r.Context(), marketProviderId, currencyMappingId, body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Provider) GetBestPrice(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["sub"])))
}
