package main

import (
	_ "fmt"
	_ "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kameike/simple_server/controller"
	_ "github.com/kameike/simple_server/model"
	_ "github.com/kameike/simple_server/model/post"
	_ "github.com/kameike/simple_server/model/session"
	_ "github.com/kameike/simple_server/util/constants"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "gopkg.in/gorp.v1"
	"log"
)

func main() {
	// Echo instance
	e := echo.New()
	// session.Create()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// e.GET("/", controller.RenderHealth)
	e.GET("/", controller.RenderHello)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
