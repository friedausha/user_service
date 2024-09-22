package mysql

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	migrationDbName = "mysql"
	mysqldriver     = "mysql"

	// DefaultTestDsn is the default url for testing postgresql in the postgres test suites
	DefaultTestDsn = "model=model password=password dbname=postgres host=localhost port=54320 sslmode=disable read_timeout=300000 write_timeout=300000"
)

// Suite struct for MySQL Suite
type Suite struct {
	suite.Suite
	DSN                     string
	DBConn                  *sql.DB
	Migration               *migration
	MigrationLocationFolder string
	DBName                  string
}

// SetupSuite setup at the beginning of test
func (s *Suite) SetupSuite() {
	viper.AutomaticEnv()
	var err error
	s.DBConn, err = sql.Open(mysqldriver, s.DSN)
	s.Require().NoError(err)
	err = s.DBConn.Ping()
	s.Require().NoError(err)
	s.Migration, err = runMigration(s.DBConn, s.MigrationLocationFolder)
	s.Require().NoError(err)
}

// TearDownSuite teardown at the end of test
func (s *Suite) TearDownSuite() {
	err := s.DBConn.Close()
	s.Require().NoError(err)
}
