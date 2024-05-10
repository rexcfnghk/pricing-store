package providercurrencyconfig

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

func providerCurrencyConfigKey(providerId int, currencyPairId int) string {
	return fmt.Sprintf("providercurrencyconfigs:%d:%d", providerId, currencyPairId)
}

func (r *RedisRepo) GetById(ctx context.Context, providerId int, currencyPairId int) (model.ProviderCurrencyConfig, error) {
	key := providerCurrencyConfigKey(providerId, currencyPairId)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.ProviderCurrencyConfig{}, errors.New("providerCurrencyConfig does not exist")
	} else if err != nil {
		return model.ProviderCurrencyConfig{}, fmt.Errorf("get providerCurrencyConfig: %w", err)
	}

	var providerCurrencyConfig model.ProviderCurrencyConfig
	err = json.Unmarshal([]byte(value), &providerCurrencyConfig)
	if err != nil {
		return model.ProviderCurrencyConfig{}, fmt.Errorf("failed to decode providerCurrencyConfig json: %w", err)
	}

	return providerCurrencyConfig, nil
}

func (r *RedisRepo) UpdateById(ctx context.Context, providerId int, currencyPairId int, config model.ProviderCurrencyConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to encode providerCurrencyConfig: %w", err)
	}

	key := providerCurrencyConfigKey(providerId, currencyPairId)

	res := r.Client.Set(ctx, key, data, 0)
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to set: %w", err)
	}

	return nil
}
