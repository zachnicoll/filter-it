package util

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// Constructs Redis client with host supplied in AWS_REDIS_ADDRESS env var
func ConnectToRedis() *redis.Client {
	redisHost := os.Getenv("AWS_REDIS_ADDRESS")
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", redisHost),
	})
}

/*
Store the given `value` against the given `key` in Redis.
Also indexes the `key` against each of the supplied filters, and
their corresponding index Sets.
*/
func CacheJSONString(ctx context.Context, key string, value string, redisClient *redis.Client) {
	_, err := redisClient.Set(ctx, key, value, 0).Result()
	if err != nil {
		fmt.Printf("Failed to fetch cache document in Redis: %v", err.Error())
	}
}

/*
	Invalidates the Redis entry for the given key, as well as the "ALL" key, "_"
*/
func InvalidateCache(ctx context.Context, key string, redisClient *redis.Client) {
	// Invalidate Redis cache entry as document was added
	redisClient.Del(ctx, key)

	// Invalidate the "ALL" entry in Redis as well
	redisClient.Del(ctx, "_")
}
