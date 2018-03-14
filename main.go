package main

import (
	"bytes"
	"fmt"
	//"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
)

var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var db_host = os.Getenv("DB_HOST")
var db_port = os.Getenv("DB_PORT")
var time = 0

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

// Handler
func hello(c echo.Context) error {
	time = time + 1
	var buf = bytes.NewBufferString("[status]\n")
	fmt.Fprintln(buf, "Hello, World!")
	fmt.Fprintf(buf, "TotalInMemoryReloadCount => %d\n", time)
	fmt.Fprintf(buf, "redis_host => %s\n", redis_host)
	fmt.Fprintf(buf, "redis_port => %s\n", redis_port)
	fmt.Fprintf(buf, "db_host => %s\n", db_host)
	fmt.Fprintf(buf, "db_port => %s\n", db_port)

	return c.String(http.StatusOK, buf.String())
}
