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
	artDB     *model.ArtDB
	accountDB *model.AccountDB
}

func (h *ArtsHandler) Init(artDB *model.ArtDB, accountDB *model.AccountDB) error {
	h.artDB = artDB
	h.accountDB = accountDB

	return nil
}

func (h ArtsHandler) PostArt(w http.ResponseWriter, r *http.Request, account dto.AccountDto) {
	w.Header().Set("Content-Type", "application/json")

	if art, err := dto.DecodeArt(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		art.AuthorId = account.Id
		if art, err := h.artDB.CreateArt(*art); err != nil {
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

func (h ArtsHandler) GetArt(w http.ResponseWriter, r *http.Request, id uint) {
	w.Header().Set("Content-Type", "application/json")

	if arts, err := h.artDB.GetArt(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(arts)
	}
}

func (h ArtsHandler) PutArt(w http.ResponseWriter, r *http.Request, art *dto.ArtDto, account dto.AccountDto) {
	w.Header().Set("Content-Type", "application/json")

	art.AuthorId = account.Id
	if art, err := h.artDB.UpdateArt(*art); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(art)
	}
}

func (h ArtsHandler) DeleteArt(w http.ResponseWriter, r *http.Request, id uint) {
	w.Header().Set("Content-Type", "application/json")
	if art, err := h.artDB.DeleteArt(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(art)
	}
}

func (h ArtsHandler) AccountAuth(w http.ResponseWriter, r *http.Request, f func(dto.AccountDto)) {
	if username, password, ok := r.BasicAuth(); !ok {
		http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
	} else if account, err := h.accountDB.GetAccountByUsername(username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if password != account.Password {
		http.Error(w, "password is incorrect", http.StatusUnauthorized)
	} else {
		f(*account)
	}
}

func (h ArtsHandler) AuthorAuth(w http.ResponseWriter, r *http.Request, id uint, f func(dto.AccountDto, dto.ArtDto)) {
	h.AccountAuth(w, r, func(account dto.AccountDto) {

		if art, err := h.artDB.GetArt(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if art.AuthorId != account.Id {
			http.Error(w, fmt.Sprintf("art #%d does not belong to '%s'", art.Id, account.Username), http.StatusUnauthorized)
		} else {
			f(account, *art)
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
		h.GetArt(w, r, uint(id))
	case http.MethodPut:
		h.AuthorAuth(w, r, uint(id), func(account dto.AccountDto, _ dto.ArtDto) {
			if art, err := dto.DecodeArt(r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			} else {
				art.Id = uint(id)
				h.PutArt(w, r, art, account)
			}
		})
	case http.MethodDelete:
		h.AuthorAuth(w, r, uint(id), func(account dto.AccountDto, _ dto.ArtDto) {
			h.DeleteArt(w, r, uint(id))
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
