package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func initDB() (*sqlx.DB, *sql.DB, error) {

	dsn := mysql.Config{
		User:                 appConfig.DB.User,
		Passwd:               appConfig.DB.Pass,
		Net:                  "tcp",
		Addr:                 appConfig.DB.Host,
		DBName:               appConfig.DB.Name,
		Timeout:              2 * time.Second,
		ReadTimeout:          2 * time.Second,
		WriteTimeout:         2 * time.Second,
		AllowNativePasswords: true,
		ParseTime:            true,

		Params: map[string]string{
			"tls": "true",
		},
	}

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {

		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return sqlx.NewDb(db, "mysql"), db, nil

}
