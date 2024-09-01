package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewSQLxDb(driverName, dbSource string) *sqlx.DB {
	conn, err := sqlx.Connect(driverName, dbSource)
	if err != nil {
		panic(err)
	}
	return conn
}
