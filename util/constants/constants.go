package constants

import (
	"fmt"
	"os"
)

var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")

var db_host = os.Getenv("DB_HOST")
var db_port = os.Getenv("DB_PORT")
var db_passowrd = os.Getenv("DB_PASSWORD")
var db_name = os.Getenv("DB_NAME")
var db_username = os.Getenv("DB_USERNAME")
var db_target = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_username, db_passowrd, db_host, db_port, db_name)

func RedisHost() string { return redis_host }
func RedisPort() string { return redis_port }

func DbHost() string     { return db_host }
func DbPort() string     { return db_port }
func DbPassword() string { return db_passowrd }
func DbName() string     { return db_name }
func DbUsername() string { return db_username }
func DbTarget() string   { return db_target }
