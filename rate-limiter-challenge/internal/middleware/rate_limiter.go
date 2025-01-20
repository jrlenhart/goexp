package middleware

import (
	"fmt"
	"github.com/JuniorLenhart/goexp/rate-limiter-challenge/internal/limiter"
	"log"
	"net/http"
)

const MessageStatusTooManyRequests = "You have reached the maximum number of requests or actions allowed within a certain time frame."

func getKey(r *http.Request) string {
	token := r.Header.Get("API_KEY")
	if token != "" {
		return fmt.Sprintf("token:%s", token)
	}
	return fmt.Sprintf("ip:%s", r.RemoteAddr)
}

func RateLimiterMiddleware(rateLimiter *limiter.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := getKey(r)
		if err := rateLimiter.AllowRequest(r.Context(), key); err != nil {
			log.Println(err)
			http.Error(w, MessageStatusTooManyRequests, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
