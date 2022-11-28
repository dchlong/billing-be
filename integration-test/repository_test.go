package integration_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"

	"github.com/dchlong/billing-be/internal/config"
	"github.com/dchlong/billing-be/internal/models"
	"github.com/dchlong/billing-be/internal/repository"
	"github.com/dchlong/billing-be/pkg/infra"
	"github.com/dchlong/billing-be/pkg/logger"
)

type RepositoryTestSuite struct {
	suite.Suite
	repo repository.IRepo
}

func (s *RepositoryTestSuite) SetupSuite() {
	appConfig, err := config.ProviderAppConfig()
	if err != nil {
		panic(err)
	}

	databaseConfig := config.ProviderDatabaseConfig(appConfig)
	iLogger, _, err := logger.ProvideLogger()
	if err != nil {
		panic(err)
	}

	db, err := infra.ProvideGormDatabase(databaseConfig, iLogger)
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	migrateTool, _, err := infra.ProvideMySQLMigrateTool(sqlDB, &infra.MigrationConfig{
		SourceFile: "file://../migrations/sql",
	})
	if err != nil {
		panic(err)
	}

	err = migrateTool.Migrate()
	if err != nil {
		panic(err)
	}

	s.repo = repository.NewSQLRepo(db)
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestBillRepo() {
	ctx := context.Background()
	callHistory := &models.CallHistory{
		UserName: "hlong",
		Duration: 1000,
	}

	err := s.repo.Bill().Create(ctx, callHistory)
	s.Require().NoError(err)
	s.Require().Greater(callHistory.ID, int64(0))

	newCallHistory, err := s.repo.Bill().FindByID(ctx, callHistory.ID)
	s.Require().NoError(err)
	s.Require().Equal(newCallHistory.ID, callHistory.ID)

	callHistories, err := s.repo.Bill().FindBy(ctx, models.NewCallHistoryFilter().WithUserName(callHistory.UserName))
	s.Require().NoError(err)
	s.Require().NotNil(callHistories)
	s.Require().Greater(len(callHistories), 0)
}
