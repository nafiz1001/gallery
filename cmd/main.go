package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ErrorJson struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *ErrorJson) Error() string {
	return e.Message
}

type Art struct {
	Id       int    `json:"id"`
	PickUp   bool   `json:"pickup"`
	Picture  string `json:"picture"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Ship     bool   `json:"ship"`
	Title    string `json:"title"`
}

func (a *Art) Update(art Art) {
	a.PickUp = art.PickUp
	a.Picture = art.Picture
	a.Price = art.Price
	a.Quantity = art.Quantity
	a.Ship = art.Ship
	a.Title = art.Title
}

func (a *Art) DecodeArt(r io.ReadCloser) error {
	var err error

	if err := json.NewDecoder(r).Decode(a); err != nil {
		return &ErrorJson{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	}

	return err
}

func HandleError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err)

		var errorJson *ErrorJson
		if !errors.As(err, &errorJson) {
			errorJson = &ErrorJson{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
		}

		w.WriteHeader(errorJson.StatusCode)
		json.NewEncoder(w).Encode(errorJson)
	}
}

func main() {
	arts := []Art{}
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		var art Art

		w.Header().Set("Content-Type", "application/json")

		if err := art.DecodeArt(r.Body); err == nil {
			art.Id = len(arts)
			arts = append(arts, art)
			json.NewEncoder(w).Encode(art)
		}

		HandleError(err, w)
	}).Methods(http.MethodPost)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(arts)
	}).Methods(http.MethodGet)

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		var err error
		var art Art

		w.Header().Set("Content-Type", "application/json")

		if err := art.DecodeArt(r.Body); err == nil {
			vars := mux.Vars(r)
			id_string, ok := vars["id"]

			if !ok {
				err = &ErrorJson{
					StatusCode: http.StatusNotFound,
					Message:    "id in path missing",
				}
			} else {
				if id, err := strconv.Atoi(id_string); err != nil {
					err = &ErrorJson{
						StatusCode: http.StatusNotFound,
						Message:    fmt.Sprintf("could not find art with id %s", id_string),
					}
				} else {
					for i := range arts {
						if arts[i].Id == id {
							arts[i].Update(art)
							json.NewEncoder(w).Encode(arts[i])
							return
						}
					}
				}
			}
		}

		HandleError(err, w)
	}).Methods(http.MethodPut)

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		id_string, ok := vars["id"]

		if !ok {
			err = &ErrorJson{
				StatusCode: http.StatusNotFound,
				Message:    "id in path missing",
			}
		} else {
			if id, err := strconv.Atoi(id_string); err != nil {
				err = &ErrorJson{
					StatusCode: http.StatusNotFound,
					Message:    fmt.Sprintf("could not find art with id %s", id_string),
				}
			} else {
				index := -1
				for i := range arts {
					if arts[i].Id == id {
						index = i
						break
					}
				}
				if index < 0 {
					err = &ErrorJson{
						StatusCode: http.StatusNotFound,
						Message:    fmt.Sprintf("could not find art with id %s", id_string),
					}
				} else {
					json.NewEncoder(w).Encode(arts[index])
					arts = append(arts[:index], arts[index+1:]...)
				}
			}
		}

		HandleError(err, w)
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
