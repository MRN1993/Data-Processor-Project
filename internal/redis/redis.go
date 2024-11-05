package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)





func InitRedis(RedisHost,RedisPort string) (*redis.Client, error){

	var ctx = context.Background() 

	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisHost + ":" + RedisPort,                     
		DB:       0,                             
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")

	return rdb, err
}