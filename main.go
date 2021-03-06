// Main
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/routers"
	"github.com/garyburd/redigo/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/urfave/negroni"
)

func main() {
	common.StartUp()

	// Set up MySQL
	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		dbPort = "3306"
	}
	dbName, ok := os.LookupEnv("DB_DBNAME")
	if !ok {
		dbName = "uhero_db_dev"
	}
	mysqlConfig := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      net.JoinHostPort(os.Getenv("DB_HOST"), dbPort),
		Loc:       time.Local,
		ParseTime: true,
		AllowNativePasswords: true,
		DBName:    dbName,
	}
	connectionString := mysqlConfig.FormatDSN()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot login to MySQL server - check all DB_* environment variables")
	}

	// Set up the FooRepository
	uhRepo := &data.FooRepository{
		DB: db,
		PortalView: "portal_v",
		SeriesView: "series_v",
		DataPointsView: "public_data_points",
		ShowHidden: false,
	}
	if view, ok := os.LookupEnv("API_PORTAL_VIEW"); ok {
		uhRepo.PortalView = view
	}
	if view, ok := os.LookupEnv("API_SERIES_VIEW"); ok {
		uhRepo.SeriesView = view
	}
	if view, ok := os.LookupEnv("API_DATAPOINTS_VIEW"); ok {
		uhRepo.DataPointsView = view
	}
	if show, ok := os.LookupEnv("API_SHOW_HIDDEN"); ok && show == "true" {
		uhRepo.ShowHidden = true
	}
	uhRepo.InitializeFoo()

	// Set up Redis
	var redisServer, authpw string
	if redisUrl, ok := os.LookupEnv("REDIS_URL"); ok {
		if u, err := url.Parse(redisUrl); err == nil {
			redisServer = u.Host // includes port where specified
			authpw, _ = u.User.Password()
		}
	}
	if redisServer == "" {
		log.Print("Valid REDIS_URL var not found; using redis @ localhost:6379")
		redisServer = "localhost:6379"
	}
	pool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisServer)
			if err != nil {
				log.Printf("*** Cannot contact redis server at %s. No caching!", redisServer)
				return nil, err
			}
			if authpw != "" {
				if _, err = c.Do("AUTH", authpw); err != nil {
					c.Close()
					log.Print("*** Redis authentication failure. No caching!")
					return nil, err
				}
			}
			log.Printf("Redis connection to %s established", redisServer)
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	ttlMinutes := 10
	cacheTtl, ok := os.LookupEnv("API_CACHE_TTL_MIN")
	if ok {
		ttlMinutes, err = strconv.Atoi(cacheTtl)
		if err != nil {
			log.Printf("*** ERROR in API_CACHE_TTL_MIN env var")
			ttlMinutes = 10
		}
	}
	log.Printf("Cache TTL is %d minutes", ttlMinutes)

	applicationRepository := uhRepo
	categoryRepository := uhRepo
	seriesRepository := uhRepo
	measurementRepository := uhRepo
	geographyRepository := uhRepo
	searchRepository := &data.SearchRepository{Categories: categoryRepository, Series: seriesRepository}
	cacheRepository := &data.CacheRepository{Pool: pool, TTL: 60 * ttlMinutes} // TTL stored in seconds

	// Get the mux router object
	router := routers.InitRoutes(
		applicationRepository,
		categoryRepository,
		seriesRepository,
		searchRepository,
		measurementRepository,
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
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
