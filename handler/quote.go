package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/rexcfnghk/pricing-store/repository/currencymapping"
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/quote"
	"github.com/shopspring/decimal"
)

type Quote struct {
	QuoteRepo           *quote.RedisRepo
	ProviderRepo        *provider.RedisRepo
	CurrencyMappingRepo *currencymapping.RedisRepo
}

type QuoteBodyModel struct {
	Base        string          `json:"base"`
	Quote       string          `json:"quote"`
	BidPrice    decimal.Decimal `json:"bidPrice"`
	BidQuantity decimal.Decimal `json:"bidQuantity"`
	AskPrice    decimal.Decimal `json:"askPrice"`
	AskQuantity decimal.Decimal `json:"askQuantity"`
	Timestamp   time.Time       `json:"timestamp"`
}

func (h *Quote) Create(w http.ResponseWriter, r *http.Request) {
	var body []QuoteBodyModel
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	var quotes []model.MarketQuote
	var errs []error
	for _, b := range body {
		quote, err := h.mapToQuote(r.Context(), marketProviderId, b)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		quotes = append(quotes, quote)
	}
	if len(errs) > 0 {
		fmt.Println("failed to map some quotes into curreny pairs: ", errs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("failed to map some quotes into curreny pairs"))
		return
	}

	errs = h.QuoteRepo.Insert(r.Context(), quotes)
	if len(errs) > 0 {
		fmt.Println("failed to insert some elemments: ", errs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("failed to insert some elemments"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Quote) mapToQuote(ctx context.Context, marketProviderId int, body QuoteBodyModel) (model.MarketQuote, error) {
	currencyMappingId, err := h.CurrencyMappingRepo.GetByCurrencyPairId(ctx, body.Base, body.Quote)
	if err != nil {
		return model.MarketQuote{}, fmt.Errorf("get currency pair: %w", err)
	}

	quote := model.MarketQuote{
		BidPrice:         body.BidPrice,
		BidQuantity:      body.BidQuantity,
		AskPrice:         body.AskPrice,
		AskQuantity:      body.AskQuantity,
		CurrencyPairId:   currencyMappingId,
		Timestamp:        body.Timestamp,
		MarketProviderId: marketProviderId,
	}

	return quote, nil
}
