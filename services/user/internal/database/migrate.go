package database

import (
	"errors"
	"log"
	"pkg/env"

	"github.com/golang-migrate/migrate/v4"
)

func RunMigrations(cmd string, version int) error {
	log.Printf("Running migration: %s", cmd)
	connStr := env.GetEnv("MIGRATION_DB_URL")
	if connStr == "" {
		return errors.New("MIGRATION_DB_URL is required")
	}

	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		return err
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
		log.Print("Applied all pending migrations")
	case "down":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			return err
		}
		log.Println("Rolled back last migration")
		return nil
	case "force":
		if version == 0 {
			log.Fatal("Need version for force command")
			return err
		}
		if err := m.Force(version); err != nil && err != migrate.ErrNoChange {
			return err
		}
		log.Printf("Forced version to %d", version)
	case "version":
		v, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			return err
		}
		log.Printf("Current version: %d, dirty: %v", v, dirty)
		return nil
	default:
		log.Fatal("Unknown commad")
	}
	return nil
}
