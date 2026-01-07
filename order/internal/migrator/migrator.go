package migrator

import (
	"database/sql"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db           *sql.DB
	migrationDir string
}

func NewMigrator(db *sql.DB, migrationDir string) *Migrator {
	return &Migrator{
		db:           db,
		migrationDir: migrationDir,
	}
}

func (m *Migrator) Up() error {
	goose.SetLogger(log.New(os.Stdout, "goose: ", log.LstdFlags)) // todo подключить zap logger
	err := goose.Up(m.db, m.migrationDir)
	if err != nil {
		return err
	}

	return nil
}
