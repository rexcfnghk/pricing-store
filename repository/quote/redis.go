package quote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	var errors []error
	for _, quote := range quotes {
		data, err := json.Marshal(quote)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to encode quote, %w", err))
			continue
		}

		key := quoteIdKey(quote.CurrencyPairId)

		res := r.Client.SAdd(ctx, key, string(data))
		if err := res.Err(); err != nil {
			errors = append(errors, fmt.Errorf("failed to add to quote set: %w", err))
			continue
		}
	}

	return errors
}

func (r *RedisRepo) GetAllByCurrencyPairId(ctx context.Context, currencyPairId int) ([]model.MarketQuote, error) {
	key := quoteIdKey(currencyPairId)

	value, err := r.Client.SMembers(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("quotes do not exist")
	} else if err != nil {
		return nil, fmt.Errorf("get quotes: %w", err)
	}

	var quotes []model.MarketQuote
	err = json.Unmarshal([]byte(strings.Join(value, "")), &quotes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode currency pair json: %w", err)
	}

	return quotes, nil
}
