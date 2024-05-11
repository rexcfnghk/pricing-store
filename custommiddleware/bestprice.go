package custommiddleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rexcfnghk/pricing-store/model"
	"github.com/rexcfnghk/pricing-store/service"
)

type BestPrice struct {
	BestPriceService *service.BestPriceService
}

func (p *BestPrice) LogBestPrice(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println("LogBestPrice")

		_, claims, _ := jwtauth.FromContext(r.Context())
		customerId, err := strconv.Atoi(fmt.Sprintf("%s", claims["sub"]))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		base, quote := r.URL.Query().Get("base"), r.URL.Query().Get("quote")
		if base == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if quote == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		currencyPair := &model.CurrencyPair{
			Base:  base,
			Quote: quote,
		}

		_, err = p.BestPriceService.GetBestPrice(r.Context(), currencyPair, customerId)

		// TODO: Add response
	})
}
