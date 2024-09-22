package mysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"strings"

	_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
)

type migration struct {
	Migrate *migrate.Migrate
}

func (m *migration) Up() (bool, error) {
	err := m.Migrate.Up()
	if err != nil {
		if errors.Is(migrate.ErrNoChange, err) {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

func (m *migration) Down() (bool, error) {
	err := m.Migrate.Down()
	if err != nil {
		if errors.Is(migrate.ErrNoChange, err) {
			return true, nil
		}
		return false, err
	}
	return true, err
}

func runMigration(dbConn *sql.DB, migrationsFolderLocation string) (*migration, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := _mysql.WithInstance(dbConn, &_mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, "mysql", driver)
	if err != nil {
		return nil, err
	}
	return &migration{Migrate: m}, nil
}
