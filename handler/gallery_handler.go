package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nafiz1001/gallery-go/arts"
	"github.com/nafiz1001/gallery-go/db"
	"github.com/nafiz1001/gallery-go/error_json"
)

func DecodeArt(r io.ReadCloser) (*arts.Art, *error_json.ErrorJson) {
	var art arts.Art

	if err := json.NewDecoder(r).Decode(&art); err != nil {
		return nil, &error_json.ErrorJson{
			Message:    err.Error(),
			StatusCode: http.StatusUnprocessableEntity,
		}
	} else {
		return &art, nil
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

type GalleryHandler struct {
	DB *db.DB
}

func (g GalleryHandler) PostArt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if art, err := DecodeArt(r.Body); err == nil {
		if art, err := g.DB.StoreArt(*art); err == nil {
			json.NewEncoder(w).Encode(art)
		} else {
			HandleError(err, w)
		}
	} else {
		HandleError(err, w)
	}
}

func (g GalleryHandler) GetArts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := g.DB.RetrieveArts(); err == nil {
		json.NewEncoder(w).Encode(arts)
	} else {
		HandleError(err, w)
	}
}

func (g GalleryHandler) GetArt(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := g.DB.RetrieveArt(id); err == nil {
		json.NewEncoder(w).Encode(arts)
	} else {
		HandleError(err, w)
	}
}

func (g GalleryHandler) PutArt(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	if art, err := DecodeArt(r.Body); err == nil {
		art.Id = id

		if art, err := g.DB.UpdateArt(*art); err == nil {
			json.NewEncoder(w).Encode(art)
		} else {
			HandleError(err, w)
		}
	} else {
		HandleError(err, w)
	}
}

func (g GalleryHandler) DeleteArt(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	if art, err := g.DB.DeleteArt(id); err == nil {
		json.NewEncoder(w).Encode(art)
	} else {
		HandleError(err, w)
	}
}

func (g GalleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	regexs := map[string]*regexp.Regexp{
		"/arts":      regexp.MustCompile("^/arts/*$"),
		"/arts/{id}": regexp.MustCompile("^/arts/([0-9]+)/*$"),
	}

	handlers := map[string](func(w http.ResponseWriter, r *http.Request)){
		"/arts": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				g.PostArt(w, r)
			case http.MethodGet:
				g.GetArts(w, r)
			default:
				http.Error(w, fmt.Sprintf("%s method not supported for %s", r.Method, r.RequestURI), http.StatusMethodNotAllowed)
			}
		},
		"/arts/{id}": func(w http.ResponseWriter, r *http.Request) {
			match := regexs["/arts/{id}"].FindStringSubmatch(r.RequestURI)
			if len(match) > 0 {
				id, _ := strconv.ParseInt(match[1], 10, 32)
				switch r.Method {
				case http.MethodGet:
					g.GetArt(w, r, int(id))
				case http.MethodPut:
					g.PutArt(w, r, int(id))
				case http.MethodDelete:
					g.DeleteArt(w, r, int(id))
				default:
					http.Error(w, fmt.Sprintf("%s method not supported for %s", r.Method, r.RequestURI), http.StatusMethodNotAllowed)
				}
			} else {
				http.Error(w, "expected uri format /arts/{id} where id is a positive integer", http.StatusBadRequest)
			}
		},
	}

	for route := range regexs {
		if regexs[route].MatchString(r.RequestURI) {
			h, ok := handlers[route]

			if ok {
				h(w, r)
				return
			} else {
				log.Printf("could not handle route '%s'", route)
			}
		}
	}

	http.NotFound(w, r)
}
