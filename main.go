package main

import (
	"database/sql"
	"fmt"
	"github.com/UHERO/rest-api/common"
	"github.com/UHERO/rest-api/data"
	"github.com/UHERO/rest-api/routers"
	"github.com/codegangsta/negroni"
	"github.com/go-sql-driver/mysql"
	"log"
	"net"
	"net/http"
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

	applicationRepository := &data.ApplicationRepository{DB: db}
	categoryRepository := &data.CategoryRepository{DB: db}
	seriesRepository := &data.SeriesRepository{DB: db}
	geographyRepository := &data.GeographyRepository{DB: db}
	feedbackRepository := &data.FeedbackRepository{DB: db}

	// Get the mux router object
	router := routers.InitRoutes(
		applicationRepository,
		categoryRepository,
		seriesRepository,
		geographyRepository,
		feedbackRepository,
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
