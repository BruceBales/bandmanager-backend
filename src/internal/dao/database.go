package dao

import (
	"database/sql"
	"fmt"

	"github.com/brucebales/bandmanager-backend/src/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewMysql() (*sql.DB, error) {
	conf := config.GetConfig()

	db, err := sql.Open(conf.MysqlDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/prim", conf.MysqlUser, conf.MysqlPass, conf.MysqlHost, conf.MysqlPort))
	if err != nil {
		return nil, err
	}
	return db, nil
}
