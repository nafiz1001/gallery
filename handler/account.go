package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	if account, err := dto.DecodeAccount(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else if acc, err := h.db.CreateAccount(*account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(acc)
	}
}

func (h AccountsHandler) GetAccountById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 32)

	if account, err := h.db.GetAccountById(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(account)
	}
}

func (h AccountsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", h.PostAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts/", h.PostAccount).Methods(http.MethodPost)

	router.HandleFunc("/accounts/{id:[0-9]+}", h.GetAccountById).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{id:[0-9]+}/", h.GetAccountById).Methods(http.MethodGet)

	router.ServeHTTP(w, r)
}
