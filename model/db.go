package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kameike/simple_server/util"
	"github.com/kameike/simple_server/util/constants"
	"gopkg.in/gorp.v1"
)

func DbMap() *gorp.DbMap {
	db, err := sql.Open("mysql", constants.DbTarget())
	util.CheckError(err, "sql.Open failed")

	return &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"MyISAM", "UTF8"}}
}
