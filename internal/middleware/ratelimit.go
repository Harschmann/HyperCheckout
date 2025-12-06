package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				// If we can't split it (maybe it's just an IP), use as is
				ip = r.RemoteAddr
			}

			userID := ip

			// Debug: Print the IP so we know it's working
			// log.Printf("Request from IP: %s", userID)

			limit := 5
			window := 1 * time.Second
			key := fmt.Sprintf("rate_limit:%s", userID)

			count, err := rdb.Incr(context.Background(), key).Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				rdb.Expire(context.Background(), key, window)
			}

			if count > int64(limit) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("429 - You are buying too fast!"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
