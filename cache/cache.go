package cache

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	stderrors "errors"

	"github.com/redis/go-redis/v9"
)

// Cache wraps a function with Redis caching

func Cache[T any](ctx context.Context, redisClient redis.UniversalClient, key string, f func() (T, error)) (T, error) {
	var model T
	val, err := redisClient.Get(ctx, key).Result()

	if err == redis.Nil || val == "" {
		// Key doesn't exist in cache
		var err error
		model, err = f()
		if err != nil {
			return model, err
		}

		// Cache the result if it's not a "not found" error
		if !stderrors.Is(err, gorm.ErrRecordNotFound) {
			modelJSON, err := json.Marshal(model)
			if err != nil {
				return model, err
			}
			if err := redisClient.Set(ctx, key, modelJSON, time.Minute).Err(); err != nil {
				return model, err
			}
		}
		return model, nil
	} else if err != nil {
		return model, err
	}

	// Unmarshal the cached result and return it
	if err := json.Unmarshal([]byte(val), &model); err != nil {
		return model, err
	}
	return model, nil
}
