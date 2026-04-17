package migrate

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // IMPORTANT
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Also required for file-based migrations
	_ "github.com/lib/pq"
)

func RunMigrations(databaseURL string) error {
	migrationsPath := "file://migrations"

	migrateDB, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func(migrateDB *migrate.Migrate) {
		err, _ := migrateDB.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close migrate instance: %w", err)
		}
	}(migrateDB)

	err = migrateDB.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if !errors.Is(err, migrate.ErrNoChange) {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

func ForceVersion(databaseURL string, version int) error {
	migrationsPath := "file://migrations"

	migrateDB, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func(migrateDB *migrate.Migrate) {
		err, _ := migrateDB.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close migrate instance: %w", err)
		}
	}(migrateDB)

	err = migrateDB.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	log.Printf("Forced migration version to %d", version)

	return nil
}
