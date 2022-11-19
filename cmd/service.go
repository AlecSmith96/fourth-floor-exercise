//go:build wireinject
// +build wireinject

package main

import (
	"github.com/AlecSmith96/fourth-floor-exercise/internal/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/entities"
	"github.com/AlecSmith96/fourth-floor-exercise/internal/rest"
	"github.com/google/wire"
)

// InitialiseService creates a new Service instance with server and logger, used by wire to generate the
// dependency injection file.
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
