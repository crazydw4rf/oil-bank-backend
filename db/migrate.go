package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Run database migrations up, down, or drop all tables.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Long:  `Apply all pending database migrations to upgrade the schema.`,
	RunE:  runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Long:  `Rollback the most recently applied migration.`,
	RunE:  runMigrateDown,
}

var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop all tables",
	Long:  `Drop all tables from the database. WARNING: This is destructive!`,
	RunE:  runMigrateDrop,
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateDropCmd)
}

func getMigrationInstance() (*migrate.Migrate, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("environment variable DATABASE_URL not found")
	}

	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating migration instance: %w", err)
	}

	return m, nil
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes to apply.")
			return nil
		}
		return fmt.Errorf("error applying migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes to apply.")
			return nil
		}
		return fmt.Errorf("error rolling back migration: %w", err)
	}

	fmt.Println("Migration rolled back successfully!")
	return nil
}

func runMigrateDrop(cmd *cobra.Command, args []string) error {
	m, err := getMigrationInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Drop(); err != nil {
		return fmt.Errorf("error dropping tables: %w", err)
	}

	fmt.Println("All tables dropped successfully!")
	return nil
}
