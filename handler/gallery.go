package handler

import (
	"log"
	"net/http"
	"regexp"
)

type GalleryHandler struct {
	ArtsHandler
}

func (h *GalleryHandler) Init() {
	h.ArtsHandler = ArtsHandler{}
	h.ArtsHandler.Init()
}

func (h GalleryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	regexs := map[string]*regexp.Regexp{
		"/arts": regexp.MustCompile("^/arts/*"),
	}

	handlers := map[string]http.Handler{
		"/arts": h.ArtsHandler,
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
