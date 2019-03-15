package dao

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Database struct {
	Driver string
	DS     string
}

func (d Database) New() (*sql.DB, error) {
	db, err := sql.Open(d.Driver, d.DS)
	if err != nil {
		return nil, err
	}
	return db, nil
}
