package photos

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed assets/migrations/*.sql
var migrationsFS embed.FS

func MigrationFS() (source.Driver, error) {
	d, err := iofs.New(migrationsFS, "assets/migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create migration source driver: %w", err)
	}

	return d, nil

}
