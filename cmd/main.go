package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nafiz1001/gallery-go/handler"
)

func main() {
	h := handler.GalleryHandler{}
	if err := h.Init(); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler: h,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
