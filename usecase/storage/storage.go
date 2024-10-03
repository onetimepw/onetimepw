package storage

import (
	"context"
	"fmt"
	"github.com/maypok86/otter"
	"github.com/onetimepw/onetimepw/domain"
	redisGate "github.com/onetimepw/onetimepw/gateway/redis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type Storage struct {
	nameSpace string
	redis     *redis.Client
	cache     *otter.CacheWithVariableTTL[string, string]
}

func New(config domain.Config) (*Storage, error) {
	if config.RedisAddr != "" {
		client, err := redisGate.NewClient(config.RedisAddr)
		if err == nil {
			zap.L().Info("redis connected", zap.String("addr", config.RedisAddr))

			return &Storage{
				redis:     client,
				nameSpace: config.NameSpace,
				cache:     nil,
			}, nil
		}

		// When redis is not connected, fallback to in-memory cache
	}

	cache, err := otter.MustBuilder[string, string](config.MemoryCapacity).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		WithVariableTTL().
		Build()
	if err != nil {
		return nil, err
	}

	zap.L().Info("redis not connected, using in-memory cache")

	return &Storage{
		redis:     nil,
		nameSpace: config.NameSpace,
		cache:     &cache,
	}, nil
}

func (s *Storage) Name() string {
	if s.redis != nil {
		return "redis"
	}

	return "in-memory"
}

func (s *Storage) Set(key string, value string, duration time.Duration) error {
	if s.redis != nil {
		return s.redis.Set(context.Background(), key, value, duration).Err()
	}

	set := s.cache.Set(key, value, duration)
	if !set {
		return fmt.Errorf("failed to set value")
	}

	return nil
}

func (s *Storage) Get(key string) (string, error) {
	if s.redis != nil {
		get := s.redis.Get(context.Background(), key)
		if get.Err() != nil {
			return "", get.Err()
		}

		return get.Val(), nil
	}

	value, found := s.cache.Get(key)
	if !found {
		return "", fmt.Errorf("value not found")
	}

	return value, nil
}

func (s *Storage) Del(key string) error {
	if s.redis != nil {
		del := s.redis.Del(context.Background(), key)
		if del.Err() != nil {
			return del.Err()
		}

		return nil
	}

	s.cache.Delete(key)

	return nil
}

func (s *Storage) Status() error {
	if s.redis != nil {
		conn := s.redis.Conn()
		defer conn.Close()
		statusCmd := conn.Ping(context.Background())
		if statusCmd == nil {
			return fmt.Errorf("can't run ping command on redis")
		}

		if statusCmd.Err() != nil {
			return statusCmd.Err()
		}

		_, err := statusCmd.Result()
		if err != nil {
			return err
		}
	}

	// Redis is not connected, no need to check cache, it's always ok
	return nil
}
