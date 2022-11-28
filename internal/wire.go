//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"

	"github.com/dchlong/billing-be/internal/config"
	"github.com/dchlong/billing-be/internal/delivery/http"
	"github.com/dchlong/billing-be/internal/repository"
	"github.com/dchlong/billing-be/internal/services"
	"github.com/dchlong/billing-be/pkg/infra"
	"github.com/dchlong/billing-be/pkg/logger"
)

func InitializeServer() (*Server, func(), error) {
	wire.Build(
		config.ProviderAppConfig,
		config.ProviderDatabaseConfig,
		ProvideRouter,
		logger.ProvideLogger,
		http.ProvideHandler,
		ProvideServer,
		services.ProvideBillService,
		infra.ProvideGormDatabase,
		repository.NewSQLRepo,
		wire.Bind(new(repository.IRepo), new(*repository.SQLRepo)),
	)

	return &Server{}, func() {}, nil
}

func InitializeMigrationTool() (infra.MigrateTool, func(), error) {
	wire.Build(
		config.ProviderAppConfig,
		config.ProviderDatabaseConfig,
		logger.ProvideLogger,
		infra.ProvideGormDatabase,
		infra.ProvideSQLDB,
		infra.ProvideMigrationConfig,
		infra.ProvideMySQLMigrateTool,
		wire.Bind(new(infra.MigrateTool), new(*infra.MySQLMigrateTool)),
	)

	return &infra.MySQLMigrateTool{}, func() {}, nil
}
