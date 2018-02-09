package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func GetConnection() (*sql.DB, error) {
	if db == nil {
		return createConnection()
	}
	return db, nil
}

func createConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", viper.GetString("db.user")+":"+viper.GetString("db.password")+"@/"+viper.GetString("db.name")+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
