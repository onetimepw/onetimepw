//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/onetimepw/onetimepw/domain"
	"github.com/onetimepw/onetimepw/endpoint/web"
	"github.com/onetimepw/onetimepw/usecase/api"
	"github.com/onetimepw/onetimepw/usecase/storage"
)

func InitApp(config domain.Config) (application, error) {
	wire.Build(
		// binds interfaces first
		wire.NewSet(
			// endpoint
			web.New,

			// binds interfaces
			wire.Bind(new(api.Storage), new(*storage.Storage)),
			wire.Bind(new(web.Storage), new(*storage.Storage)),
			// gateway

			// usecase
			api.New,
			storage.New,
		),
		// app
		newApplication,
	)
	return application{}, nil
}
