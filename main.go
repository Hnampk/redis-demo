package main

import (
	"context"
	"fmt"
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

	channelKey := "CHANNEL_*"
	dbChannels, err := getAllKeyValue(channelKey)
	if err != nil {
		log.Printf("error while getAllKeyValue of %s: %s\n", channelKey, err.Error())
	} else {
		log.Printf("dbChannels: %+v\n", dbChannels)
	}

}

func getAllKeyValue(searchkey string) (map[string]string, error) {
	res := make(map[string]string)
	// Scan all keys
	var cursor uint64

	ctx := context.Background()
	for {
		var keys []string
		var err error
		keys, cursor, err = redisService.Client.Scan(ctx, cursor, searchkey, 0).Result()
		if err != nil {
			fmt.Printf("Scan key: %s , error: %s\n", searchkey, err)
			return res, err
		}

		for _, key := range keys {
			val, err := redisService.Client.Get(ctx, key).Result()
			if err != nil {
				fmt.Printf("getAllKeyValue: error to get value from key  %s error %s\n", key, err)
				fmt.Println("Error to get value from key: ", key)
				continue
			}
			//Add to value
			res[key] = val
		}

		if cursor == 0 { // no more keys
			break
		}
	}
	return res, nil
}
