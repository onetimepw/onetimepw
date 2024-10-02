package redis

import (
	"app/domain"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewClient(config domain.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(config.Redis.Addr)
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
