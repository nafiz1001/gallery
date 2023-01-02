package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nafiz1001/gallery-go/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nafiz1001/gallery-go/handler"
)

func main() {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	h := handler.GalleryHandler{}
	err = h.Init(&model.DB{
		GormDB: gormDB,
	})
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler: h,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Listening to localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
