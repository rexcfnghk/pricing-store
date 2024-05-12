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
	"github.com/rexcfnghk/pricing-store/repository/currencypair"
	"github.com/rexcfnghk/pricing-store/repository/provider"
	"github.com/rexcfnghk/pricing-store/repository/quote"
	"github.com/shopspring/decimal"
)

type Quote struct {
	QuoteRepo        *quote.RedisRepo
	ProviderRepo     *provider.RedisRepo
	CurrencyPairRepo *currencypair.RedisRepo
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
		marketQuote, err := h.mapToQuote(r.Context(), marketProviderId, b)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		quotes = append(quotes, marketQuote)
	}
	if len(errs) > 0 {
		fmt.Println("failed to map some quotes into currency pairs: ", errs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("failed to map some quotes into currency pairs"))
		return
	}

	errs = h.QuoteRepo.Insert(r.Context(), quotes)
	if len(errs) > 0 {
		fmt.Println("failed to insert some elements: ", errs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("failed to insert some elements"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Quote) mapToQuote(ctx context.Context, marketProviderId int, body QuoteBodyModel) (model.MarketQuote, error) {
	currencyPairId, err := h.CurrencyPairRepo.GetByCurrencyPairId(ctx, body.Base, body.Quote)
	if err != nil {
		return model.MarketQuote{}, fmt.Errorf("get currency pair: %w", err)
	}

	marketQuote := model.MarketQuote{
		BidPrice:         body.BidPrice,
		BidQuantity:      body.BidQuantity,
		AskPrice:         body.AskPrice,
		AskQuantity:      body.AskQuantity,
		CurrencyPairId:   currencyPairId,
		Timestamp:        body.Timestamp,
		MarketProviderId: marketProviderId,
	}

	return marketQuote, nil
}
