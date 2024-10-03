package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewClient(redisAddr string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	ping := client.Ping(context.Background())

	_, err = ping.Result()
	if err != nil {
		return nil, fmt.Errorf("redisClient: %s", err.Error())
	}

	return client, nil
}
