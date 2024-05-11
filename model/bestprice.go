package model

import "github.com/shopspring/decimal"

type BestPrice struct {
	BidPrice         decimal.Decimal `json:"bidPrice"`
	BidQuantity      decimal.Decimal `json:"bidQuantity"`
	AskPrice         decimal.Decimal `json:"askPrice"`
	AskQuantity      decimal.Decimal `json:"askQuantity"`
	MarketProviderId int             `json:"marketProviderId"`
}
