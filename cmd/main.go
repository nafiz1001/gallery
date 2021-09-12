package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nafiz1001/gallery-go/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// connection string

	// postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
	// psqlconn, ok := os.LookupEnv("DATABASE_URL")
	// if !ok {
	// 	log.Fatal("expected env var DATABASE_URL")
	// }

	psqlconn := "./database.db"

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	h := handler.GalleryHandler{}
	err = h.Init(db)
	CheckError(err)

	srv := &http.Server{
		Handler: h,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
