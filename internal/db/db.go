package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", os.Getenv("DB_CONN_STRING"))

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}
