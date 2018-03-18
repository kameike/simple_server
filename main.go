package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/gorp.v1"
	"log"
	"net/http"
	"os"
	"time"
)

var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")

var db_host = os.Getenv("DB_HOST")
var db_port = os.Getenv("DB_PORT")
var db_passowrd = os.Getenv("DB_PASSWORD")
var db_name = os.Getenv("DB_NAME")
var db_username = os.Getenv("DB_USERNAME")
var db_target = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_username, db_passowrd, db_host, db_port, db_name)

var time_count = 0
var redisClient = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", redis_host, redis_port),
	Password: "", // no password set
	DB:       0,  // use default DB
})

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
	time_count = time_count + 1
	var buf = bytes.NewBufferString("[env]\n")
	fmt.Fprintln(buf, "Hello, World!")
	fmt.Fprintf(buf, "TotalInMemoryReloadCount => %d\n", time_count)
	fmt.Fprintf(buf, "redis_host => %s\n", redis_host)
	fmt.Fprintf(buf, "redis_port => %s\n", redis_port)
	fmt.Fprintf(buf, "db_host => %s\n", db_host)
	fmt.Fprintf(buf, "db_port => %s\n", db_port)
	fmt.Fprintf(buf, "db_passowrd => %s\n", db_passowrd)
	fmt.Fprint(buf, "\n")

	// redis helth check
	fmt.Fprintln(buf, "[redis_state]")
	_, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Fprintf(buf, "error => %s\n", err)
	} else {
		fmt.Fprintf(buf, "redis available\n")
	}
	fmt.Fprintf(buf, "target => %s:%s\n", redis_host, redis_port)
	fmt.Fprint(buf, "\n")

	// mysql helth check
	var db, _ = sql.Open("mysql", db_target)
	fmt.Fprintln(buf, "[mysql_state]")
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(buf, "error => %s\n", err)
	} else if db != nil {
		fmt.Fprintf(buf, "database available\n")
	}
	fmt.Fprintf(buf, "target => %s\n", db_target)
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
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	// create two posts
	p1 := newPost("Go 1.1 released!", "Lorem ipsum lorem ipsum")
	p2 := newPost("Go 1.2 released!", "Lorem ipsum lorem ipsum")

	// insert rows - auto increment PKs will be set properly after the insert
	err = dbmap.Insert(&p1, &p2)
	checkErr(err, "Insert failed")

	// use convenience SelectInt
	count, err := dbmap.SelectInt("select count(*) from posts")
	checkErr(err, "select count(*) failed")
	log.Println("Rows after inserting:", count)

	// update a row
	p2.Title = "Go 1.2 is better than ever"
	count, err = dbmap.Update(&p2)
	checkErr(err, "Update failed")
	log.Println("Rows updated:", count)

	// fetch one row - note use of "post_id" instead of "Id" since column is aliased
	//
	// Postgres users should use $1 instead of ? placeholders
	// See 'Known Issues' below
	//
	err = dbmap.SelectOne(&p2, "select * from posts where post_id=?", p2.Id)
	checkErr(err, "SelectOne failed")
	log.Println("p2 row:", p2)

	// fetch all rows
	var posts []Post
	_, err = dbmap.Select(&posts, "select * from posts order by post_id")
	checkErr(err, "Select failed")
	log.Println("All rows:")
	for x, p := range posts {
		log.Printf("    %d: %v\n", x, p)
	}

	// delete row by PK
	count, err = dbmap.Delete(&p1)
	checkErr(err, "Delete failed")
	log.Println("Rows deleted:", count)

	// delete row manually via Exec
	_, err = dbmap.Exec("delete from posts where post_id=?", p2.Id)
	checkErr(err, "Exec failed")

	// confirm count is zero
	count, err = dbmap.SelectInt("select count(*) from posts")
	checkErr(err, "select count(*) failed")
	log.Println("Row count - should be zero:", count)

	log.Println("Done!")
}

type Post struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      int64 `db:"post_id"`
	Created int64
	Title   string `db:",size:50"`               // Column size set to 50
	Body    string `db:"article_body,size:1024"` // Set both column name and size
}

func newPost(title, body string) Post {
	return Post{
		Created: time.Now().UnixNano(),
		Title:   title,
		Body:    body,
	}
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("mysql", db_target)

	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"MyISAM", "UTF8"}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
