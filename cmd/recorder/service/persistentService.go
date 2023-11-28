package service

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/int2xx9/daikin-airconditioner/cmd/recorder/config"
	_ "github.com/lib/pq"
)

type PersistentService struct {
	config config.Configuration
}

func NewPersistentService(config config.Configuration) PersistentService {
	return PersistentService{
		config: config,
	}
}

func (s PersistentService) Migrate() error {
	m, err := migrate.New("file://migrations", s.config.Database.ConnectionString())
	if err != nil {
		return err
	}
	m.Log = &migrateLogger{}
	m.Up()
	return nil
}

type migrateLogger struct{}

func (migrateLogger) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

func (migrateLogger) Verbose() bool {
	return true
}
