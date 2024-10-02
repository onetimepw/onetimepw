package api

import (
	"app/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type API struct {
	config domain.Config
	redis  *redis.Client
}

func New(config domain.Config, redis *redis.Client) *API {
	return &API{
		config: config,
		redis:  redis,
	}
}

func (a *API) redisKey(id string) string {
	redisKey := fmt.Sprintf("%s:secrets:%s", a.config.Redis.NameSpace, id)

	return redisKey
}

func (a *API) Create(text string, password string, duration time.Duration) (id string, resultPassword string, err error) {
	if text == "" {
		return "", "", fmt.Errorf("text is empty")
	}

	maxLen := 5000
	if len(text) > maxLen {
		return "", "", fmt.Errorf("text is too long, max length is %d bytes", maxLen)
	}

	if password == "" {
		pwd, err := randomPassword(6)
		if err != nil {
			return "", "", err
		}
		password = pwd
	}

	encryptData, err := encrypt(text, password)
	if err != nil {
		return "", "", err
	}

	base64Encoded := toBase64([]byte(encryptData))

	id = uuid.NewString()

	key := a.redisKey(id)

	err = a.redis.Set(context.Background(), key, base64Encoded, duration).Err()
	if err != nil {
		return "", "", err
	}

	return id, password, nil
}

func (a *API) Has(id string) bool {
	redisKey := a.redisKey(id)

	get := a.redis.Get(context.Background(), redisKey)

	return get.Err() == nil
}

func (a *API) Get(id, password string) (string, error) {
	redisKey := a.redisKey(id)

	keyData, err := a.redis.Get(context.Background(), redisKey).Result()
	if err != nil {
		return "", errors.New("not found")
	}

	base64Decoded, err := fromBase64(keyData)
	if err != nil {
		return "", errors.New("invalid key")
	}

	decryptData, err := decrypt(string(base64Decoded), password)
	if err != nil {
		// cipher: message authentication failed
		return "", errors.New("invalid password")
	}

	// delete key after read
	err = a.redis.Del(context.Background(), redisKey).Err()
	if err != nil {
		zap.L().Error("failed to delete key", zap.String("key", redisKey), zap.Error(err))
	}

	return decryptData, nil
}
