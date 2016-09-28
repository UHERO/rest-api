package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uhero/rest-api/common"
	"github.com/uhero/rest-api/data"
	"github.com/uhero/rest-api/routers"
	"log"
	"net/http"
	"os"
)

func main() {
	common.StartUp()
	// Set up MySQL
	connectionString := fmt.Sprintf(
		"%s:%s@%s(%s)/%s?parseTime=true&loc=US%%2FHawaii",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		common.AppConfig.Protocol,
		common.AppConfig.Server,
		common.AppConfig.Database,
	)
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

	// Get the mux router object
	router := routers.InitRoutes(applicationRepository, categoryRepository)
	// Create a negroni instance
	n := negroni.Classic()
	n.UseHandler(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}
	log.Printf("Listening on %s...", server.Addr)
	server.ListenAndServe()
}
