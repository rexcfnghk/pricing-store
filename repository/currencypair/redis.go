package currencypair

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

type CurrencyPairId = int

func currencyPairKey(base string, quote string) string {
	return fmt.Sprintf("currencypairs:%s:%s", base, quote)
}

func (r *RedisRepo) GetByCurrencyPairId(ctx context.Context, base string, quote string) (CurrencyPairId, error) {
	key := currencyPairKey(base, quote)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, errors.New("currency pair does not exist")
	} else if err != nil {
		return 0, fmt.Errorf("get currency pair: %w", err)
	}

	var currenyPairId CurrencyPairId
	err = json.Unmarshal([]byte(value), &currenyPairId)
	if err != nil {
		return 0, fmt.Errorf("failed to decode currency pair json: %w", err)
	}

	return currenyPairId, nil
}
