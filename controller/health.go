package controller

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kameike/simple_server/util/constants"
	"github.com/labstack/echo"
	"net/http"
)

var time_count = 0

// Handler
func RenderHealth(c echo.Context) error {
	time_count = time_count + 1
	var buf = bytes.NewBufferString("[env]\n")
	fmt.Fprintln(buf, "Hello, World!")
	fmt.Fprintf(buf, "TotalInMemoryReloadCount => %d\n", time_count)
	fmt.Fprintf(buf, "redis_host => %s\n", constants.RedisHost())
	fmt.Fprintf(buf, "redis_port => %s\n", constants.RedisPort())
	fmt.Fprintf(buf, "db_host => %s\n", constants.DbHost())
	fmt.Fprintf(buf, "db_port => %s\n", constants.DbPort())
	fmt.Fprintf(buf, "db_passowrd => %s\n", constants.DbPassword())
	fmt.Fprint(buf, "\n")

	// redis helth check
	fmt.Fprintln(buf, "[redis state]")
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", constants.RedisHost(), constants.RedisPort()),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Fprintf(buf, "error => %s\n", err)
	} else {
		fmt.Fprintf(buf, "redis available\n")
	}
	fmt.Fprint(buf, "\n")

	// mysql helth check
	fmt.Fprintln(buf, "[mysql state]")
	var db, _ = sql.Open("mysql", constants.DbTarget())
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(buf, "error => %s\n", err)
	} else if db != nil {
		fmt.Fprintf(buf, "database available\n")
	}
	fmt.Fprint(buf, "\n")

	return c.String(http.StatusOK, buf.String())
}
