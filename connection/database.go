package connection

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func DbConnect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	err = runMigrations(databaseURL)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database migrations: No new changes applied.")
		} else {

			return nil, fmt.Errorf("failed to run database migrations: %w", err)
		}
	} else {
		log.Println("Database migrations ran successfully.")
	}

	return db, nil
}

func runMigrations(connectionString string) error {
	m, err := migrate.New("file://./migrations/migrations", connectionString)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		var dirtyErr migrate.ErrDirty
		if errors.As(err, &dirtyErr) {
			log.Printf("Dirty database version %d found. Forcing to this version to clear the dirty state.", dirtyErr.Version)
			if forceErr := m.Force(dirtyErr.Version); forceErr != nil {
				return fmt.Errorf("failed to force migration version: %w", forceErr)
			}
			// After forcing, try to migrate up again.
			log.Println("Retrying migrations after forcing version.")
			err = m.Up()
		}
	}

	// After all attempts, return the final error state.
	// The caller will handle ErrNoChange.
	return err
}
