package quote

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rexcfnghk/pricing-store/model"
)

type RedisRepo struct {
	Client *redis.Client
}

func quoteIdKey(quote model.MarketQuote) string {
	return fmt.Sprintf("quotes:%d", quote.CurrencyPairId)
}

func (r *RedisRepo) Insert(ctx context.Context, quotes []model.MarketQuote) []error {
	var errors []error
	for _, quote := range quotes {
		data, err := json.Marshal(quote)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to encode quote, %w", err))
			continue
		}

		key := quoteIdKey(quote)

		res := r.Client.SAdd(ctx, key, string(data))
		if err := res.Err(); err != nil {
			errors = append(errors, fmt.Errorf("failed to add to quote set: %w", err))
			continue
		}
	}

	return errors
}
