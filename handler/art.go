package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/nafiz1001/gallery-go/dto"
	"github.com/nafiz1001/gallery-go/model"
)

type ArtsHandler struct {
	db *model.ArtDB
}

func (h *ArtsHandler) Init() error {
	h.db = &model.ArtDB{}
	return h.db.Init()
}

func (h ArtsHandler) PostArt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var art dto.ArtDto

	if err := json.NewDecoder(r.Body).Decode(&art); err == nil {
		if art, err := h.db.StoreArt(art); err == nil {
			json.NewEncoder(w).Encode(art)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h ArtsHandler) GetArts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := h.db.RetrieveArts(); err == nil {
		json.NewEncoder(w).Encode(arts)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h ArtsHandler) GetArt(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := h.db.RetrieveArt(id); err == nil {
		json.NewEncoder(w).Encode(arts)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h ArtsHandler) PutArt(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	var art dto.ArtDto

	if err := json.NewDecoder(r.Body).Decode(&art); err == nil {
		art.Id = id
		if art, err := h.db.UpdateArt(art); err == nil {
			json.NewEncoder(w).Encode(art)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h ArtsHandler) DeleteArt(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	if art, err := h.db.DeleteArt(id); err == nil {
		json.NewEncoder(w).Encode(art)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h ArtsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	regexs := map[string]*regexp.Regexp{
		"/arts":      regexp.MustCompile("^/arts/*$"),
		"/arts/{id}": regexp.MustCompile("^/arts/([0-9a-z]+)/*$"),
	}

	handlers := map[string](func(w http.ResponseWriter, r *http.Request)){
		"/arts": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				h.PostArt(w, r)
			case http.MethodGet:
				h.GetArts(w, r)
			default:
				http.Error(w, fmt.Sprintf("%s method not supported for %s", r.Method, r.RequestURI), http.StatusMethodNotAllowed)
			}
		},
		"/arts/{id}": func(w http.ResponseWriter, r *http.Request) {
			match := regexs["/arts/{id}"].FindStringSubmatch(r.RequestURI)
			if len(match) > 0 {
				switch r.Method {
				case http.MethodGet:
					h.GetArt(w, r, match[1])
				case http.MethodPut:
					h.PutArt(w, r, match[1])
				case http.MethodDelete:
					h.DeleteArt(w, r, match[1])
				default:
					http.Error(w, fmt.Sprintf("%s method not supported for %s", r.Method, r.RequestURI), http.StatusMethodNotAllowed)
				}
			} else {
				http.Error(w, "expected uri format /arts/{id} where {id} is a positive integer", http.StatusBadRequest)
			}
		},
	}

	for route := range regexs {
		if regexs[route].MatchString(r.RequestURI) {
			handler, ok := handlers[route]

			if ok {
				handler(w, r)
				return
			} else {
				log.Printf("could not handle route '%s'", route)
			}
		}
	}

	http.NotFound(w, r)
}
