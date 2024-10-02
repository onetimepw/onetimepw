//go:build wireinject
// +build wireinject

package main

import (
	"app/domain"
	"app/endpoint/web"
	"app/gateway/redis"
	"app/usecase/api"
	"github.com/google/wire"
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
