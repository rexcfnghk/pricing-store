package customer

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

func customerIdKey(customerId int) string {
	return fmt.Sprintf("customers:%d", customerId)
}

func (r *RedisRepo) GetById(ctx context.Context, id int) (model.Customer, error) {
	key := customerIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Customer{}, errors.New("customer does not exist")
	} else if err != nil {
		return model.Customer{}, fmt.Errorf("get customer: %w", err)
	}

	var customer model.Customer
	err = json.Unmarshal([]byte(value), &customer)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to decode customer json: %w", err)
	}

	return customer, nil
}
