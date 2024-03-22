package sqlite

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed internal/migrations
var migrations embed.FS

func (a *AppendOnlyStore) MigrateDB() (err error) {
	db, err := a.db.DB()
	if err != nil {
		return fmt.Errorf("unable to retrieve database connection: %w", err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("unable to create migration driver: %w", err)
	}

	fs, err := iofs.New(migrations, "internal/migrations")
	if err != nil {
		return fmt.Errorf("unable to create migration fs: %w", err)
	}

	instance, err := migrate.NewWithInstance("iofs", fs, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration instance: %w", err)
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to migrate: %w", err)
	}

	return nil
}
