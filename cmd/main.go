package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nafiz1001/gallery-go/arts"
	"github.com/nafiz1001/gallery-go/db"
	"github.com/nafiz1001/gallery-go/error_json"
)

func DecodeArt(r io.ReadCloser) (*arts.Art, error) {
	var art *arts.Art

	if err := json.NewDecoder(r).Decode(&art); err != nil {
		return nil, &error_json.ErrorJson{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	} else {
		return art, err
	}
}

func HandleError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err)

		var errorJson *error_json.ErrorJson
		if !errors.As(err, &errorJson) {
			errorJson = &error_json.ErrorJson{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
		}

		http.Error(w, errorJson.Message, errorJson.StatusCode)
	}
}

func main() {
	db := db.Init()
	r := mux.NewRouter()

	r.HandleFunc("/arts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if art, err := DecodeArt(r.Body); err == nil {
			if art, err := db.StoreArt(*art); err == nil {
				json.NewEncoder(w).Encode(art)
			} else {
				HandleError(err, w)
			}
		} else {
			HandleError(err, w)
		}
	}).Methods(http.MethodPost)

	r.HandleFunc("/arts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if arts, err := db.RetrieveAllArt(); err == nil {
			json.NewEncoder(w).Encode(arts)
		} else {
			HandleError(err, w)
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/arts/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if art, err := DecodeArt(r.Body); err == nil {
			vars := mux.Vars(r)
			id_string := vars["id"]
			id, _ := strconv.Atoi(id_string)
			art.Id = id

			if art, err = db.UpdateArt(*art); err == nil {
				json.NewEncoder(w).Encode(art)
			} else {
				HandleError(err, w)
			}
		} else {
			HandleError(err, w)
		}
	}).Methods(http.MethodPut)

	r.HandleFunc("/arts/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		id_string := vars["id"]
		id, _ := strconv.Atoi(id_string)

		if art, err := db.DeleteArt(id); err == nil {
			json.NewEncoder(w).Encode(art)
		} else {
			HandleError(err, w)
		}
	}).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
