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
	artDB          *model.ArtDB
	accountDB      *model.AccountDB
	accountsArtsDB *model.AccountsArtsDB
}

func (h *ArtsHandler) Init(artDB *model.ArtDB, accountDB *model.AccountDB, accountsArtsDB *model.AccountsArtsDB) error {
	h.artDB = artDB
	h.accountDB = accountDB
	h.accountsArtsDB = accountsArtsDB

	return nil
}

func (h ArtsHandler) PostArt(w http.ResponseWriter, r *http.Request, account dto.AccountDto) {
	w.Header().Set("Content-Type", "application/json")

	var art dto.ArtDto

	if err := json.NewDecoder(r.Body).Decode(&art); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		if art, err := h.accountsArtsDB.AddArt(account, art); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(art)
		}
	}
}

func (h ArtsHandler) GetArts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := h.artDB.GetArts(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(arts)
	}
}

func (h ArtsHandler) GetArt(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := h.artDB.GetArt(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(arts)
	}
}

func (h ArtsHandler) PutArt(w http.ResponseWriter, r *http.Request, art *dto.ArtDto) {
	w.Header().Set("Content-Type", "application/json")

	if art, err := h.artDB.UpdateArt(*art); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(art)
	}
}

func (h ArtsHandler) DeleteArt(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")
	if art, err := h.accountsArtsDB.DeleteArt(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(art)
	}
}

func (h ArtsHandler) AccountAuth(w http.ResponseWriter, r *http.Request, f func(dto.AccountDto)) {
	if username, password, ok := r.BasicAuth(); !ok {
		http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
	} else {
		if account, err := h.accountDB.GetAccount(username); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if password != account.Password {
				http.Error(w, "password is incorrect", http.StatusUnauthorized)
			} else {
				f(*account)
			}
		}
	}
}

func (h ArtsHandler) AuthorAuth(w http.ResponseWriter, r *http.Request, id string, f func(dto.AccountDto)) {
	h.AccountAuth(w, r, func(account dto.AccountDto) {
		if !h.accountsArtsDB.IsAuthor(account, id) {
			http.Error(w, fmt.Sprintf("art does not belong to '%s'", account.Username), http.StatusUnauthorized)
		} else {
			f(account)
		}
	})
}

func (h ArtsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	regexs := map[string]*regexp.Regexp{
		"/arts":      regexp.MustCompile("^/arts/*$"),
		"/arts/{id}": regexp.MustCompile("^/arts/([^/]+)/*$"),
	}

	handlers := map[string](func(w http.ResponseWriter, r *http.Request)){
		"/arts": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				h.AccountAuth(w, r, func(account dto.AccountDto) {
					h.PostArt(w, r, account)
				})
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
					h.AuthorAuth(w, r, match[1], func(account dto.AccountDto) {
						var art dto.ArtDto
						if err := json.NewDecoder(r.Body).Decode(&art); err != nil {
							http.Error(w, err.Error(), http.StatusUnprocessableEntity)
						} else {
							art.Id = match[1]
							h.PutArt(w, r, &art)
						}
					})
				case http.MethodDelete:
					h.AuthorAuth(w, r, match[1], func(account dto.AccountDto) {
						h.DeleteArt(w, r, match[1])
					})
				default:
					http.Error(w, fmt.Sprintf("%s method not supported for %s", r.Method, r.RequestURI), http.StatusMethodNotAllowed)
				}
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