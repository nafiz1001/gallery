package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nafiz1001/gallery-go/dto"
	"github.com/nafiz1001/gallery-go/model"
)

type ArtsHandler struct {
	artDB         *model.ArtDB
	accountDB     *model.AccountDB
	accountArtsDB *model.AccountArtsDB
}

func (h *ArtsHandler) Init(artDB *model.ArtDB, accountDB *model.AccountDB, accountArtsDB *model.AccountArtsDB) error {
	h.artDB = artDB
	h.accountDB = accountDB
	h.accountArtsDB = accountArtsDB

	return nil
}

func (h ArtsHandler) PostArt(w http.ResponseWriter, r *http.Request, account dto.AccountDto) {
	w.Header().Set("Content-Type", "application/json")

	var art dto.ArtDto

	if err := json.NewDecoder(r.Body).Decode(&art); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		if art, err := h.accountArtsDB.AddArt(account, art); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(art)
		}
	}
}

func (h ArtsHandler) GetArts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := h.accountArtsDB.GetArts(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(arts)
	}
}

func (h ArtsHandler) GetArt(w http.ResponseWriter, r *http.Request, id int) {
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

func (h ArtsHandler) DeleteArt(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	if art, err := h.accountArtsDB.DeleteArt(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(art)
	}
}

func (h ArtsHandler) AccountAuth(w http.ResponseWriter, r *http.Request, f func(dto.AccountDto)) {
	if username, password, ok := r.BasicAuth(); !ok {
		http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
	} else {
		if account, err := h.accountDB.GetAccountByUsername(username); err != nil {
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

func (h ArtsHandler) AuthorAuth(w http.ResponseWriter, r *http.Request, id int, f func(dto.AccountDto)) {
	h.AccountAuth(w, r, func(account dto.AccountDto) {
		if !h.accountArtsDB.IsAuthor(account, id) {
			http.Error(w, fmt.Sprintf("art does not belong to '%s'", account.Username), http.StatusUnauthorized)
		} else {
			f(account)
		}
	})
}

func (h ArtsHandler) ArtsFuncHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.AccountAuth(w, r, func(account dto.AccountDto) {
			h.PostArt(w, r, account)
		})
	case http.MethodGet:
		h.GetArts(w, r)
	}
}

func (h ArtsHandler) ArtByIdFuncHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 32)

	switch r.Method {
	case http.MethodGet:
		h.GetArt(w, r, int(id))
	case http.MethodPut:
		h.AuthorAuth(w, r, int(id), func(account dto.AccountDto) {
			var art dto.ArtDto
			if err := json.NewDecoder(r.Body).Decode(&art); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			} else {
				art.Id = int(id)
				h.PutArt(w, r, &art)
			}
		})
	case http.MethodDelete:
		h.AuthorAuth(w, r, int(id), func(account dto.AccountDto) {
			h.DeleteArt(w, r, int(id))
		})
	}
}

func (h ArtsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/arts", h.ArtsFuncHandler).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/arts/", h.ArtsFuncHandler).Methods(http.MethodPost, http.MethodGet)

	router.HandleFunc("/arts/{id:[0-9]+}", h.ArtByIdFuncHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	router.HandleFunc("/arts/{id:[0-9]+}/", h.ArtByIdFuncHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	router.ServeHTTP(w, r)
}
