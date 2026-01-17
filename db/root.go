package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management CLI for Oil Bank",
	Long:  `A command-line interface for managing database migrations and seeding data for the Oil Bank application.`,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}
