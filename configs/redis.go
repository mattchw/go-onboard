package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
)

// Client instance
var RDB *redis.Client = ConnectRedis()

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     EnvRedisAddress(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// ping redis
	res, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Redis", res)
	return rdb
}
