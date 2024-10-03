//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/onetimepw/onetimepw/domain"
	"github.com/onetimepw/onetimepw/endpoint/web"
	"github.com/onetimepw/onetimepw/gateway/redis"
	"github.com/onetimepw/onetimepw/usecase/api"
)

func InitApp(config domain.Config) (application, error) {
	wire.Build(
		// binds interfaces first
		wire.NewSet(
			// endpoint
			web.New,

			// binds interfaces
			// gateway
			redis.NewClient,

			// usecase
			api.New,
		),
		// app
		newApplication,
	)
	return application{}, nil
}
