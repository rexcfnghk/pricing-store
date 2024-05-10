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
	CurrencyPair     CurrencyPair
	MarketProviderId int32
}

type CurrencyPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}
