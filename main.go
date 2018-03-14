package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strconv"
)

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

var time = 0

// Handler
func hello(c echo.Context) error {
	time = time + 1
	return c.String(http.StatusOK, "Hello, World!"+strconv.Itoa(time))
}
