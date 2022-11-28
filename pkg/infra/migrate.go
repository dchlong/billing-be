package infra

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

type MigrateTool interface {
	Migrate() error
}

type MySQLMigrateTool struct {
	migrate *migrate.Migrate
}

func (mt *MySQLMigrateTool) Migrate() error {
	version, dirty, err := mt.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if dirty {
		err = mt.migrate.Force(int(version) - 1)
		if err != nil {
			return err
		}
	}

	err = mt.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

type MigrationConfig struct {
	SourceFile string
}

func ProvideMigrationConfig() *MigrationConfig {
	return &MigrationConfig{
		SourceFile: "file://migrations/sql",
	}
}

func ProvideMySQLMigrateTool(db *sql.DB, cfg *MigrationConfig) (*MySQLMigrateTool, func(), error) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, func() {}, err
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.SourceFile, "mysql", driver)
	if err != nil {
		return nil, func() {}, err
	}

	return &MySQLMigrateTool{
			migrate: m,
		}, func() {
			_ = db.Close()
		}, nil
}
