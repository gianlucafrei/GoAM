package sqlite_adapter

import (
	"database/sql"
	"embed"
	"fmt"
	"goiam/internal/logger"
	"strings"

	_ "modernc.org/sqlite" // sqlite driver
)

type Config struct {
	Driver string // "sqlite" or "pgx"
	DSN    string
}

func Init(cfg Config) (*sql.DB, error) {
	var err error
	database, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return database, nil
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Uses embed fs to load migrations and run them over the database connection in order
func RunMigrations(db *sql.DB) error {

	// Open migrations folder
	migrations, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	// go over all files in migrations folder
	for _, migration := range migrations {
		// read file
		migrationFile, err := migrationsFS.ReadFile("migrations/" + migration.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration: %w", err)
		}

		// Only expecute up migrations
		if strings.Contains(migration.Name(), "up.sql") {
			// Log name of migration
			logger.DebugNoContext("Running migration: %s", migration.Name())

			// run migration
			_, err = db.Exec(string(migrationFile))
			if err != nil {
				return fmt.Errorf("failed to run migration: %w", err)
			}
		}
	}

	return nil
}
