package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/rexcfnghk/pricing-store/repository/quote"
	"github.com/shopspring/decimal"
)

type Quote struct {
	Repo *quote.RedisRepo
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

	var quotes []model.MarketQuote
	for _, b := range body {
		quotes = append(quotes, mapToQuote(marketProviderId, b))
	}

	errs := h.Repo.Insert(r.Context(), quotes)
	if len(errs) > 0 {
		fmt.Println("failed to insert some elemments: ", errs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("failed to insert some elemments"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func mapToQuote(marketProviderId int, body QuoteBodyModel) model.MarketQuote {
	quote := model.MarketQuote{
		BidPrice:    body.BidPrice,
		BidQuantity: body.BidQuantity,
		AskPrice:    body.AskPrice,
		AskQuantity: body.AskQuantity,
		CurrencyPair: model.CurrencyPair{
			Base:  body.Base,
			Quote: body.Quote,
		},
		Timestamp:        body.Timestamp,
		MarketProviderId: marketProviderId,
	}

	return quote
}
