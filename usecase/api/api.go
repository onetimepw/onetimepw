package api

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/onetimepw/onetimepw/domain"
	"go.uber.org/zap"
	"time"
)

type Storage interface {
	Set(key string, value string, duration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type API struct {
	config  domain.Config
	storage Storage
}

func New(config domain.Config,
	store Storage) *API {
	return &API{
		config:  config,
		storage: store,
	}
}

func (a *API) idToKey(id string) string {
	redisKey := fmt.Sprintf("%s:secrets:%s", a.config.NameSpace, id)

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

	key := a.idToKey(id)

	err = a.storage.Set(key, base64Encoded, duration)
	if err != nil {
		return "", "", err
	}

	return id, password, nil
}

func (a *API) Has(id string) bool {
	key := a.idToKey(id)
	_, err := a.storage.Get(key)

	return err == nil
}

func (a *API) Get(id, password string) (string, error) {
	key := a.idToKey(id)

	text, err := a.storage.Get(key)
	if err != nil {
		return "", err
	}

	base64Decoded, err := fromBase64(text)
	if err != nil {
		return "", errors.New("invalid key")
	}

	decryptData, err := decrypt(string(base64Decoded), password)
	if err != nil {
		// cipher: message authentication failed
		return "", errors.New("invalid password")
	}

	// delete key after read
	err = a.storage.Del(key)
	if err != nil {
		zap.L().Error("failed to delete key", zap.String("key", key), zap.Error(err))
	}

	return decryptData, nil
}
