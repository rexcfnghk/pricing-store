package quote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rexcfnghk/pricing-store/model"
)

type RedisRepo struct {
	Client *redis.Client
}

func quoteIdKey(currencyPairId int) string {
	return fmt.Sprintf("quotes:%d", currencyPairId)
}

func (r *RedisRepo) Insert(ctx context.Context, quotes []model.MarketQuote) []error {
	var errs []error
	for _, quote := range quotes {
		data, err := json.Marshal(quote)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to encode quote, %w", err))
			continue
		}

		key := quoteIdKey(quote.CurrencyPairId)

		res := r.Client.SAdd(ctx, key, string(data))
		if err := res.Err(); err != nil {
			errs = append(errs, fmt.Errorf("failed to add to quote set: %w", err))
			continue
		}
	}

	return errs
}

func (r *RedisRepo) GetAllByCurrencyPairId(ctx context.Context, currencyPairId int) ([]model.MarketQuote, error) {
	key := quoteIdKey(currencyPairId)

	values, err := r.Client.SMembers(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get quotes: %w", err)
	}

	if noQuotes := len(values) == 0; noQuotes {
		return nil, nil
	}

	var quotes []model.MarketQuote
	for _, value := range values {
		quote, err := mapToMarketQuote(value)
		if err != nil {
			return nil, fmt.Errorf("map json to quote failed: %w", err)
		}

		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func mapToMarketQuote(value string) (model.MarketQuote, error) {
	var quote model.MarketQuote
	err := json.Unmarshal([]byte(value), &quote)
	if err != nil {
		return model.MarketQuote{}, fmt.Errorf("failed to decode market quote json: %w", err)
	}

	return quote, nil
}
