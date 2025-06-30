package connection

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

//var DB *sql.DB

func DbConnect() *sql.DB {
	// Connect to PostgreSQL database
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		connectionString = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		fmt.Println("Warning: DATABASE_URL environment variable not set. Using default local PostgreSQL connection string.")
	}
	var err error
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	err = runMigrations(connectionString)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Database migrations: No new changes applied.")
		} else {

			panic(fmt.Sprintf("Failed to run database migrations: %v", err))
		}
	} else {

		fmt.Println("Database migrations ran successfully.")
	}
	return db

}

func runMigrations(connectionString string) error {
	m, err := migrate.New("file://./migrations/migrations", connectionString)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
