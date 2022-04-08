package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	username string = os.Getenv("DATABASE_USERNAME")
	password string = os.Getenv("DATABASE_PASSWORD")
	database string = os.Getenv("DATABASE_NAME")
)

func GetConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v?parseTime=true", username, password, database))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db, nil
}
