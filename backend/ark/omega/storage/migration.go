package storage

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrationVer uint = 20240212020000

type MigrateAction func(m *migrate.Migrate) error

func RunMigration(driver database.Driver, action MigrateAction) error {
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return err
	}

	return action(m)
}

func MigrateUp() MigrateAction {
	return func(m *migrate.Migrate) error {
		return m.Up()
	}
}

func MigrateDown() MigrateAction {
	return func(m *migrate.Migrate) error {
		return m.Down()
	}
}

func VerifyMigrationVersion(expected uint) MigrateAction {
	return func(m *migrate.Migrate) error {
		actual, dirty, err := m.Version()
		if err != nil {
			return err
		}
		if dirty {
			return fmt.Errorf("dirty migrate version (version: %v)", actual)
		}
		if expected != actual {
			return fmt.Errorf("different migrate version (actual: %v, expected: %v)", actual, expected)
		}
		return nil
	}
}
