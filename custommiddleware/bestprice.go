package custommiddleware

import (
	"fmt"
	"net/http"
)

func LogBestPrice(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println("LogBestPrice")
	})
}
