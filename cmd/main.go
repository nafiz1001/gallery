package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nafiz1001/gallery-go/db"
	"github.com/nafiz1001/gallery-go/handler"
)

func main() {
	h := handler.GalleryHandler{
		DB: db.Init(),
	}

	srv := &http.Server{
		Handler: h,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
