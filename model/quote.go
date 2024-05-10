package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type MarketQuote struct {
	BidPrice         decimal.Decimal `json:"bidPrice"`
	BidQuantity      decimal.Decimal `json:"bidQuantity"`
	AskPrice         decimal.Decimal `json:"askPrice"`
	AskQuantity      decimal.Decimal `json:"askQuantity"`
	Timestamp        time.Time       `json:"timestamp"`
	MarketProviderId int             `json:"marketProviderId"`
	CurrencyPairId   int
}

type CurrencyPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}
