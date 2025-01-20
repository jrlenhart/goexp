package limiter

import (
	"context"
	"fmt"
	"log"
	"time"
)

type RateLimiter struct {
	storage         Storage
	maxRequestIP    int
	maxRequestToken int
	blockTimeIP     int
	blockTimeToken  int
}

func NewRateLimiter(storage Storage, maxRequestIP, maxRequestToken, blockTimeIP, blockTimeToken int) *RateLimiter {
	return &RateLimiter{
		storage:         storage,
		maxRequestIP:    maxRequestIP,
		maxRequestToken: maxRequestToken,
		blockTimeIP:     blockTimeIP,
		blockTimeToken:  blockTimeToken,
	}
}

func (rl *RateLimiter) AllowRequest(ctx context.Context, key string) error {
	if isBlocked, err := rl.isBlocked(ctx, key); err != nil {
		return err
	} else if isBlocked {
		return fmt.Errorf("key %s is blocked", key)
	}

	count, err := rl.incrementRequest(ctx, key)
	if err != nil {
		return err
	}
	log.Printf("key %s counter: %d", key, count)

	if count == 1 {
		if err = rl.setExpiration(ctx, key); err != nil {
			return err
		}
	}

	limit := rl.maxRequestIP
	blockTime := rl.blockTimeIP
	if isToken(key) {
		limit = rl.maxRequestToken
		blockTime = rl.blockTimeToken
	}

	if count > limit {
		if err = rl.blockKey(ctx, key, time.Duration(blockTime)*time.Second); err != nil {
			return err
		}
		return fmt.Errorf("key %s exceeded the limit and is now blocked", key)
	}

	return nil
}

func (rl *RateLimiter) incrementRequest(ctx context.Context, key string) (int, error) {
	count, err := rl.storage.Increment(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to increment request count for key %s: %w", key, err)
	}
	return count, nil
}

func (rl *RateLimiter) setExpiration(ctx context.Context, key string) error {
	if err := rl.storage.Expire(ctx, key, time.Duration(1)*time.Second); err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}

func (rl *RateLimiter) isBlocked(ctx context.Context, key string) (bool, error) {
	isBlocked, err := rl.storage.IsBlocked(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check if key %s is blocked: %w", key, err)
	}
	return isBlocked, nil
}

func (rl *RateLimiter) blockKey(ctx context.Context, key string, blockTime time.Duration) error {
	if err := rl.storage.BlockKey(ctx, key, blockTime); err != nil {
		return fmt.Errorf("failed to block key %s: %w", key, err)
	}
	return nil
}

func isToken(key string) bool {
	return len(key) > 5 && key[:6] == "token:"
}
