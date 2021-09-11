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

type AccountsHandler struct {
	db *model.AccountDB
}

func (h *AccountsHandler) Init(db *model.AccountDB) error {
	h.db = db
	return nil
}

func (h AccountsHandler) PostAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var account dto.AccountDto

	if err := json.NewDecoder(r.Body).Decode(&account); err == nil {
		if acc, err := h.db.CreateAccount(account); err == nil {
			json.NewEncoder(w).Encode(acc)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h AccountsHandler) GetAccount(w http.ResponseWriter, r *http.Request, username string) {
	w.Header().Set("Content-Type", "application/json")

	if account, err := h.db.GetAccount(username); err == nil {
		json.NewEncoder(w).Encode(account)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h AccountsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	regexs := map[string]*regexp.Regexp{
		"/accounts":      regexp.MustCompile("^/accounts/*$"),
		"/accounts/{id}": regexp.MustCompile("^/accounts/([^/]]+)/*$"),
	}

	handlers := map[string](func(w http.ResponseWriter, r *http.Request)){
		"/accounts": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				h.PostAccount(w, r)
			default:
				http.Error(w, fmt.Sprintf("%s method not supported for %s", r.Method, r.RequestURI), http.StatusMethodNotAllowed)
			}
		},
		"/accounts/{id}": func(w http.ResponseWriter, r *http.Request) {
			match := regexs["/accounts/{id}"].FindStringSubmatch(r.RequestURI)
			if len(match) > 0 {
				switch r.Method {
				case http.MethodGet:
					h.GetAccount(w, r, match[1])
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
