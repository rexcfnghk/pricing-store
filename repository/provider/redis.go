package provider

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

func providerIdKey(providerId int) string {
	return fmt.Sprintf("providers:%d", providerId)
}

func (r *RedisRepo) GetById(ctx context.Context, id int) (model.Provider, error) {
	key := providerIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Provider{}, errors.New("provider does not exist")
	} else if err != nil {
		return model.Provider{}, fmt.Errorf("get provider: %w", err)
	}

	var quote model.Provider
	err = json.Unmarshal([]byte(value), &quote)
	if err != nil {
		return model.Provider{}, fmt.Errorf("failed to decode provider json: %w", err)
	}

	return quote, nil
}
