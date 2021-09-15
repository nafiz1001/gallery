package handler

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nafiz1001/gallery-go/model"
)

type GalleryHandler struct {
	artDB           *model.ArtDB
	accountDB       *model.AccountDB
	accountArtsDB   *model.AccountArtsDB
	artsHandler     ArtsHandler
	accountsHandler AccountsHandler
}

func (h *GalleryHandler) Init(db *sql.DB) error {
	h.artDB = &model.ArtDB{}
	if err := h.artDB.Init(db); err != nil {
		return err
	}

	h.accountDB = &model.AccountDB{}
	if err := h.accountDB.Init(db); err != nil {
		return err
	}

	h.accountArtsDB = &model.AccountArtsDB{}
	if err := h.accountArtsDB.Init(db, h.accountDB, h.artDB); err != nil {
		return err
	}

	h.artsHandler = ArtsHandler{}
	if err := h.artsHandler.Init(h.artDB, h.accountDB, h.accountArtsDB); err != nil {
		return err
	}

	h.accountsHandler = AccountsHandler{}
	if err := h.accountsHandler.Init(h.accountDB); err != nil {
		return err
	}

	return nil
}

func (h GalleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.PathPrefix("/accounts").Handler(h.accountsHandler)
	router.PathPrefix("/accounts/").Handler(h.accountsHandler)

	router.PathPrefix("/arts").Handler(h.artsHandler)
	router.PathPrefix("/arts/").Handler(h.artsHandler)

	router.ServeHTTP(w, r)
}
