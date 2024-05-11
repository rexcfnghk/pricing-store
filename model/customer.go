package model

import (
	"github.com/shopspring/decimal"
)

type Customer struct {
	RatingFactor decimal.Decimal `json:"ratingFactor"`
}
