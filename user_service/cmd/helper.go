package cmd

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

// InitDBMySQL initialize postgres database
func InitDBMySQL() *sqlx.DB {
	host := MustHaveEnv("DB_HOST")
	portStr := MustHaveEnv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err, "DB_PORT is not well set ")
	}
	fmt.Println(port)
	user := MustHaveEnv("DB_USER")
	password := MustHaveEnv("DB_PASSWORD")
	dbname := MustHaveEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, host, port, dbname)
	fmt.Println(dsn)

	if err != nil {
		log.Fatalln(err)
	}
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Restrict the connection count to reduce DB load
	maxConnStr := GetOptionalEnv("MAX_DB_CONNECTION", "100")
	maxConn, err := strconv.Atoi(maxConnStr)
	if err != nil {
		log.Fatalln(errors.New("MAX_DB_CONNECTION is not well set "))
	}
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(3000)
	db.SetConnMaxIdleTime(30)

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	return db
}

// MustHaveEnv ensure the ENV exists, otherwise it will crash the app
func MustHaveEnv(key string) string {
	val, isSet := viper.GetString(key), viper.IsSet(key)
	if !isSet {
		log.Fatalf("%s is not set", key)
	}

	if val == "" {
		log.Fatalf("%s is not valid", key)
	}
	return val
}

func GetOptionalEnv(key string, defaultVal string) string {
	val, isSet := viper.GetString(key), viper.IsSet(key)
	if !isSet {
		return defaultVal
	}
	return val
}
