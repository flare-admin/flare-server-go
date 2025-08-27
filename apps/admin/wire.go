//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/flare-admin/flare-server-go/apps/admin/server"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/google/wire"
)

// wireApp init application.
func wireApp(*configs.Bootstrap, *configs.Data) (*app, func(), error) {
	panic(wire.Build(server.ProviderSet, newApp))
}
