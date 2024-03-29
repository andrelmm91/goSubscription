package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscription/data"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var webPort = "80"
var counts int64

func main() {
	log.Println("Stating subscription service")

	// connect to database
	db := initDB()

	// create sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroups
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session: session,
		DB: db,
		InfoLog: infoLog,
		ErrorLog: errorLog,
		Wait: &wg,
		Models: data.New(db),
	}

	// set up mails

	// listen for signals
	go app.listenForShutDown()

	// listen for web connection
	app.serve()
}

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

/////////////////////////////////////////////
// Connecting to DB
/////////////////////////////////////////////

func initDB() *sql.DB {
	// To Do connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	return conn
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

/////////////////////////////////////////////
// Session
/////////////////////////////////////////////

func initSession() *scs.SessionManager {
	gob.Register(data.User{}) // store user in the session
	// set up Session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func () (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return redisPool
}

func (app *Config) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
 	// Waits for a value to be sent on the quit channel.
	// When a signal is received (either SIGINT or SIGTERM), the execution continues past this line.
	<-quit 

	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	//perform any cleanup tasks
	app.InfoLog.Println("run cleanup tasks")

	// block waitgroup is empty
	app.Wait.Wait()

	app.InfoLog.Println("closing channels and shutting down application..")
}
