package handler

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"

	"github.com/nafiz1001/gallery-go/model"
)

type GalleryHandler struct {
	artDB           *model.ArtDB
	accountDB       *model.AccountDB
	accountsArtsDB  *model.AccountsArtsDB
	artsHandler     ArtsHandler
	accountsHandler AccountsHandler
}

func (h *GalleryHandler) Init(db *sql.DB) error {
	h.artDB = &model.ArtDB{}
	h.accountDB = &model.AccountDB{}
	h.accountsArtsDB = &model.AccountsArtsDB{}

	if err := h.artDB.Init(db); err != nil {
		return err
	}
	if err := h.accountDB.Init(db); err != nil {
		return err
	}
	if err := h.accountsArtsDB.Init(db, h.accountDB, h.artDB); err != nil {
		return err
	}

	h.artsHandler = ArtsHandler{}
	if err := h.artsHandler.Init(h.artDB, h.accountDB, h.accountsArtsDB); err != nil {
		return err
	}

	h.accountsHandler = AccountsHandler{}
	if err := h.accountsHandler.Init(h.accountDB); err != nil {
		return err
	}

	return nil
}

func (h GalleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	regexs := map[string]*regexp.Regexp{
		"/arts":     regexp.MustCompile("^/arts/*"),
		"/accounts": regexp.MustCompile("^/accounts/*"),
	}

	handlers := map[string]http.Handler{
		"/arts":     h.artsHandler,
		"/accounts": h.accountsHandler,
	}

	for route := range regexs {
		if regexs[route].MatchString(r.RequestURI) {
			handler, ok := handlers[route]

			if ok {
				handler.ServeHTTP(w, r)
				return
			} else {
				log.Printf("could not handle route '%s'", route)
			}
		}
	}

	http.NotFound(w, r)
}
