package handler

import (
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

func (h *GalleryHandler) Init() {
	h.artDB = &model.ArtDB{}
	h.accountDB = &model.AccountDB{}
	h.accountsArtsDB = &model.AccountsArtsDB{}

	h.artDB.Init()
	h.accountDB.Init()
	h.accountsArtsDB.Init(h.accountDB, h.artDB)

	h.artsHandler = ArtsHandler{}
	h.artsHandler.Init(h.artDB, h.accountDB, h.accountsArtsDB)

	h.accountsHandler = AccountsHandler{}
	h.accountsHandler.Init(h.accountDB)
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
