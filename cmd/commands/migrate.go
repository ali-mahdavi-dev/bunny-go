package commands

import (
	"fmt"
	"net/url"

	"errors"
	"log"

	"github.com/ali-mahdavi-dev/bunny-go/internal/account/adapter/migrations"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/spf13/cobra"
)

var ErrMigrationFileNameRequired = errors.New("migration name is required")

func dbmateDB() *dbmate.DB {
	cnn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Tehran",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.DbName, cfg.Postgres.SSLMode)
	u, err := url.Parse(cnn)
	if err != nil {
		panic(fmt.Errorf("invalid DB connection string: %w", err))
	}

	dbConn := dbmate.New(u)
	dbConn.FS = migrations.Migrations
	dbConn.MigrationsDir = []string{"./"}
	dbConn.AutoDumpSchema = false

	return dbConn
}

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "handle database migration actions",
	}

	migrateMake := &cobra.Command{
		Use:   "make",
		Short: "create new migrate",
		RunE: func(_ *cobra.Command, args []string) error {
			initializeConfigs()
			if len(args) == 0 {
				return ErrMigrationFileNameRequired
			}

			return makeMigration(args[0])
		},
	}

	migrateUp := &cobra.Command{
		Use:   "up",
		Short: "migrate the database",
		RunE: func(_ *cobra.Command, _ []string) error {
			initializeConfigs()

			return migrate()
		},
	}

	migrateDown := &cobra.Command{
		Use:   "down",
		Short: "rollback database migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			initializeConfigs()

			return migrateRollback()
		},
	}

	migrateStatus := &cobra.Command{
		Use:   "status",
		Short: "get migration status",
		RunE: func(_ *cobra.Command, _ []string) error {
			initializeConfigs()

			return migrateStatus()
		},
	}

	cmd.AddCommand(migrateMake)
	cmd.AddCommand(migrateUp)
	cmd.AddCommand(migrateDown)
	cmd.AddCommand(migrateStatus)

	return cmd
}

func migrateStatus() error {
	log.Println("Migrations:")
	migrations, err := dbmateDB().FindMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}
	for _, m := range migrations {
		if m.Applied {
			log.Println("[✅]", m.Version, m.FilePath)
		} else {
			log.Println("[❌]", m.Version, m.FilePath)
		}
	}

	return nil
}

func migrate() error {
	log.Println("Applying Migrations:")
	err := dbmateDB().CreateAndMigrate()
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func migrateRollback() error {
	err := dbmateDB().Rollback()
	if err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

func makeMigration(name string) error {
	db := dbmateDB()
	db.MigrationsDir = []string{"migrations"}

	err := db.NewMigration(name)
	if err != nil {
		return fmt.Errorf("failed to create database migration: %w", err)
	}

	log.Println("new migration created: ", name)

	return nil
}
