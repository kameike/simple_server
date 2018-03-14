package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
)

var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var db_host = os.Getenv("DB_HOST")
var db_port = os.Getenv("DB_PORT")
var db_passowrd = os.Getenv("DB_PASSWORD")
var db_name = os.Getenv("DB_NAME")
var db_username = os.Getenv("DB_USERNAME")
var time = 0
var redisClient = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", redis_host, redis_port),
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {
	// Echo instance
	test()
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
	var buf = bytes.NewBufferString("[env]\n")
	fmt.Fprintln(buf, "Hello, World!")
	fmt.Fprintf(buf, "TotalInMemoryReloadCount => %d\n", time)
	fmt.Fprintf(buf, "redis_host => %s\n", redis_host)
	fmt.Fprintf(buf, "redis_port => %s\n", redis_port)
	fmt.Fprintf(buf, "db_host => %s\n", db_host)
	fmt.Fprintf(buf, "db_port => %s\n", db_port)
	fmt.Fprintf(buf, "db_passowrd => %s\n", db_passowrd)
	fmt.Fprint(buf, "\n")

	// redis helth check
	fmt.Fprintln(buf, "[redis_state]")
	pong, err := redisClient.Ping().Result()
	fmt.Fprintln(buf, pong, err)
	fmt.Fprint(buf, "\n")

	var target = fmt.Sprintf("root:%s@tcp(%s:%s)/", db_passowrd, db_host, db_port)
	// mysql helth check
	var db, dbInitError = sql.Open("mysql", target)
	fmt.Fprintln(buf, "[mysql_state]")
	fmt.Println(target)
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(buf, "error => %s\n", dbInitError)
	}
	if db != nil {
		fmt.Fprintf(buf, "db is here\n", dbInitError)
	}
	fmt.Fprint(buf, "\n")

	return c.String(http.StatusOK, buf.String())
}

func exampleredisClient() {
	err := redisClient.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := redisClient.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func test() {
	var target = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_username, db_passowrd, db_host, db_port, db_name)
	fmt.Println(target)
	db, err := sql.Open("mysql", target)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d", squareNum)
}
