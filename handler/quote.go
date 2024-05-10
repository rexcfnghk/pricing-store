package handler

import (
	"fmt"
	"net/http"

	"github.com/rexcfnghk/pricing-store/repository/quote"
)

type Quote struct {
	Repo *quote.RedisRepo
}

func (o *Quote) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")
}
