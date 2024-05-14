package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/shopspring/decimal"
)

var ctx = context.TODO()

func TestGetById_ReturnsCustomer_GivenIdExists(t *testing.T) {
	// Arrange
	db, mock := redismock.NewClientMock()
	const customerId = 1
	redisKey := fmt.Sprintf("customers:%d", customerId)
	const ratingFactor = 1
	want := model.Customer{RatingFactor: decimal.NewFromInt(ratingFactor)}
	json, _ := json.Marshal(want)
	mock.ExpectGet(redisKey).SetVal(string(json))

	sut := &RedisRepo{Client: db}

	// Act
	actual, err := sut.GetById(ctx, customerId)

	// Assert
	if err != nil || !actual.RatingFactor.Equal(want.RatingFactor) {
		t.Errorf("Cannot find matching customer, got %v, want: %v, error: %v", *actual, want, err)
	}
}
