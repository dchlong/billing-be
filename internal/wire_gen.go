// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/dchlong/billing-be/internal/config"
	"github.com/dchlong/billing-be/internal/delivery/http"
	"github.com/dchlong/billing-be/internal/repository"
	"github.com/dchlong/billing-be/internal/services"
	"github.com/dchlong/billing-be/pkg/infra"
	"github.com/dchlong/billing-be/pkg/logger"
)

// Injectors from wire.go:

func InitializeServer() (*Server, func(), error) {
	appConfig, err := config.ProviderAppConfig()
	if err != nil {
		return nil, nil, err
	}
	iLogger, cleanup, err := logger.ProvideLogger()
	if err != nil {
		return nil, nil, err
	}
	engine := ProvideRouter(iLogger)
	databaseConfig := config.ProviderDatabaseConfig(appConfig)
	db, err := infra.ProvideGormDatabase(databaseConfig, iLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sqlRepo := repository.NewSQLRepo(db)
	billService := services.ProvideBillService(sqlRepo, appConfig)
	handler := http.ProvideHandler(billService, iLogger)
	server, cleanup2, err := ProvideServer(appConfig, engine, handler, iLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return server, func() {
		cleanup2()
		cleanup()
	}, nil
}

func InitializeMigrationTool() (infra.MigrateTool, func(), error) {
	appConfig, err := config.ProviderAppConfig()
	if err != nil {
		return nil, nil, err
	}
	databaseConfig := config.ProviderDatabaseConfig(appConfig)
	iLogger, cleanup, err := logger.ProvideLogger()
	if err != nil {
		return nil, nil, err
	}
	db, err := infra.ProvideGormDatabase(databaseConfig, iLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sqlDB, err := infra.ProvideSQLDB(db)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	migrationConfig := infra.ProvideMigrationConfig()
	mySQLMigrateTool, cleanup2, err := infra.ProvideMySQLMigrateTool(sqlDB, migrationConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return mySQLMigrateTool, func() {
		cleanup2()
		cleanup()
	}, nil
}