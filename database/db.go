package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

const (
	SSLModeDisable SSLMode = "disable"
)

type SSLMode string

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {
	connectionStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, databaseName, user, password, sslMode)

	var err error
	DB, err = sqlx.Open("postgres", connectionStr)
	if err != nil {
		return err
	}
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("database ping failed %w", err)
	}
	fmt.Println("Database connected successfully")
	return migrateUp(DB)
}

func migrateUp(db *sqlx.DB) error {
	fmt.Println("Database Migration Starting...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No New Migration")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}
	fmt.Println("Database Migration Successful")
	return nil
}

func Tx(fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	err = fn(tx)
	return
}
