package main

import (
	"fmt"
	"github.com/jrlenhart/goexp/rate-limiter-challenge/internal/configuration"
	"github.com/jrlenhart/goexp/rate-limiter-challenge/internal/limiter"
	"github.com/jrlenhart/goexp/rate-limiter-challenge/internal/middleware"
	"log"
	"net/http"
)

func main() {
	appConfig := configuration.NewApplication()
	if err := appConfig.Load(); err != nil {
		log.Fatal("failed to load application environment", err)
	}

	redisStorage := limiter.NewRedisStorage(appConfig.RedisHost, appConfig.RedisPort)
	rateLimiter := limiter.NewRateLimiter(redisStorage, appConfig.MaxRequestIP, appConfig.MaxRequestToken, appConfig.BlockTimeIP, appConfig.BlockTimeToken)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("RateLimiter Challenge!"))
	})

	rateLimiterMiddleware := middleware.RateLimiterMiddleware(rateLimiter, mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: rateLimiterMiddleware,
	}

	log.Println("server running on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server error: ", err)
	}
}
