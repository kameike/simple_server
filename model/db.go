package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kameike/simple_server/util"
	"github.com/kameike/simple_server/util/constants"
	"gopkg.in/gorp.v1"
)

var db *sql.DB

func DbMap() *gorp.DbMap {
	if db.Stats().OpenConnections < 3 {
		var err error
		db, err = sql.Open("mysql", constants.DbTarget())
		util.CheckError(err, "sql.Open failed")
	}

	return &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"MyISAM", "UTF8"}}
}
