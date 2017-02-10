// Main
package main

import (
	"database/sql"
	"fmt"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/routers"
	"github.com/codegangsta/negroni"
	"github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	common.StartUp()

	// Set up MySQL
	mysqlConfig := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      net.JoinHostPort(os.Getenv("DB_HOST"), "3306"),
		Loc:       time.Local,
		ParseTime: true,
		DBName:    "uhero_db_dev",
	}
	connectionString := mysqlConfig.FormatDSN()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Start MySQL Server!")
	}

	var redis_server, authpw string
	if redis_url, present := os.LookupEnv("REDIS_URL"); present {
		if u, ok := url.Parse(redis_url); ok == nil {
			redis_server = u.Host // includes port where specified
			authpw, _ = u.User.Password()
		}
	}
	if redis_server == "" {
		log.Printf("Valid REDIS_URL var not found; using redis @ localhost:6379")
		redis_server = "localhost:6379"
	}
	redis_conn, err := redis.Dial("tcp", redis_server)
	if err != nil {
		log.Printf("*** Cannot contact redis server at %s. No caching!", redis_server)
	} else if authpw != "" {
		if _, err = redis_conn.Do("AUTH", authpw); err != nil {
			redis_conn.Close()
			redis_conn = nil
			log.Printf("*** Redis authentication failure. No caching!")
		}
	} else {
		defer redis_conn.Close()
	}

	applicationRepository := &data.ApplicationRepository{DB: db}
	categoryRepository := &data.CategoryRepository{DB: db}
	seriesRepository := &data.SeriesRepository{DB: db}
	geographyRepository := &data.GeographyRepository{DB: db}
	cacheRepository := &data.CacheRepository{DB: redis_conn}

	// Get the mux router object
	router := routers.InitRoutes(
		applicationRepository,
		categoryRepository,
		seriesRepository,
		geographyRepository,
		cacheRepository,
	)
	// Create a negroni instance
	n := negroni.Classic()
	n.UseHandler(router)

	port := os.Getenv("GO_REST_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: n,
	}
	log.Printf("Listening on %s...", server.Addr)
	server.ListenAndServe()
}
