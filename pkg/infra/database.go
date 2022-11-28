package infra

import (
	"database/sql"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"

	"github.com/dchlong/billing-be/pkg/logger"
)

type DatabaseConfig struct {
	DataSource string `json:"data_source"`
}

func ProvideGormDatabase(cfg *DatabaseConfig, ilogger logger.ILogger) (*gorm.DB, error) {
	gormlogger := logger.NewGormLogger(ilogger)
	gormlogger.SetAsDefault()
	const maxRetries = 10
	var db *gorm.DB
	gerr := backoff.Retry(func() error {
		var err error
		db, err = gorm.Open(mysql.Open(cfg.DataSource), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 gormlogger,
		})

		return err
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries))
	if gerr != nil {
		return nil, gerr
	}

	return db, nil
}

func ProvideSQLDB(gormDB *gorm.DB) (*sql.DB, error) {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}
