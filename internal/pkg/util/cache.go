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
Create a key that is the concatenation of all filters in ascending order.
E.g. [2, 1, 3] creates the key "123"
*/
func ConstructRedisKey(filters []int) string {
	redisKey := "_" // "_" is the "root" key, where a search for ALL documents is stored
	for _, filter := range filters {
		redisKey = fmt.Sprintf("%s%d", redisKey, filter)
	}

	return redisKey
}

/*
Store `value` against `key` in Redis, and print an error if one occurs.
Does NOT panic if this operation fails.
*/
func cacheValue(ctx context.Context, key string, value string, redisClient *redis.Client) {
	_, err := redisClient.Set(ctx, key, value, 0).Result()
	if err != nil {
		fmt.Printf("Failed to fetch cache document in Redis: %v", err.Error())
	}
}

/*
Add Redis key (e.g. "_123") to each index Set.
See explanation of index Sets in Notion:
https://www.notion.so/Caching-Strategy-8b25c9a2f1354186b99e1e313e45d3ae
*/
func indexFilters(ctx context.Context, key string, filters []int, redisClient *redis.Client) {
	for _, filter := range filters {
		// Format index like "x:index", where x is a filter e.g. "3:index"
		_, err := redisClient.SAdd(ctx, fmt.Sprintf("%d:index", filter), key).Result()
		if err != nil {
			fmt.Printf("Failed to add key to Set index in Redis: %v", err.Error())
		}
	}
}

/*
Store the given `value` against the given `key` in Redis.
Also indexes the `key` against each of the supplied filters, and
their corresponding index Sets.
*/
func CacheJSONString(ctx context.Context, key string, value string, filters []int, redisClient *redis.Client) {
	// Store JSON string in Redis
	cacheValue(ctx, key, value, redisClient)

	// Store the key against each of the filter indices
	indexFilters(ctx, key, filters, redisClient)
}
