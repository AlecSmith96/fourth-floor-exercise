//go:build wireinject
// +build wireinject
package main

import (
	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/entities"
	"github.com/AlecSmith96/fourth-floor-exercise/rest"
	"github.com/google/wire"
)

// InitialiseService creates a new Service instance with server and logger
func InitialiseService() (Service, error) {
	wire.Build(
		adapters.NewLogger,
		adapters.NewConfig,
		wire.FieldsOf(new(*entities.Config), "Rest", "Logging", "Auth"),
		adapters.NewTwitchAdapter,
		adapters.NewAnalyticsAdapter,
		rest.NewRouter,
		rest.NewHTTPServer,
		wire.Struct(new(Service), "*"),
	)

	return Service{}, nil
}
