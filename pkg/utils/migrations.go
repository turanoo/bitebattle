package utils

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func RunMigrations(dbURL string, migrationsPath string, log *logrus.Entry) error {
	log.Info("Starting database migrations...")
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	var closeErr error
	defer func() {
		if err, _ := m.Close(); err != nil {
			closeErr = err
		}
	}()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close migrate instance: %w", closeErr)
	}
	log.Info("Database migrations ran successfully.")
	return nil
}
