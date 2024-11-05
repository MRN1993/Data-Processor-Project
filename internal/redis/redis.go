package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"data-processor-project/internal/logs"
	"go.uber.org/zap"
)





func InitRedis(RedisHost,RedisPort string) (*redis.Client, error){

	var ctx = context.Background() 

	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisHost + ":" + RedisPort,                     
		DB:       0,                             
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logs.Logger.Fatal("Failed to connect to Redis",zap.Error(err))
	}

	logs.Logger.Info("Connected to Redis successfully")

	return rdb, err
}