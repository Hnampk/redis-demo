package main

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v8"
)

type RedisWrapper struct {
	Client *redis.Client
}

var (
	redisService *RedisWrapper
)

func init() {
	redisMasterName := "master"
	sentinels := []string{"localhost:26379", "localhost:26379", "localhost:26379"}
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    redisMasterName,
		SentinelAddrs: sentinels,
	})

	redisService = &RedisWrapper{
		Client: rdb,
	}
}

func getKey(ctx context.Context, key string) (value string, err error) {

	return redisService.Client.Get(ctx, key).Result()
}

func main() {
	mykey := "account001"
	value, err := getKey(context.Background(), mykey)
	if err != nil {
		if err != redis.Nil {
			// we got log here, which is "context deadline exceeded"
			log.Printf("error while get key from redis: %s", err.Error())
		}

		log.Printf("key %s not found in redis", mykey)
	}

	log.Printf("key %s found in redis, value=%s", mykey, value)
}
