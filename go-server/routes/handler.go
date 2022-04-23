package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ultrasound-client/go-server/facade"
	"net/http"
	"os"
)

type Handler struct {
	Service facade.ProxyFacade
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	// Static routes
	r.Handle("/", h.ClientHandler())
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", h.ClientHandler()))
	// Proxy route
	r.PathPrefix("/api/").Handler(h.ServiceHandler())
	// Health check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			return
		}
	})
	return r
}

func (h Handler) ClientHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		// build paths with request and config file
		res, svcErr := h.Service.Client(urlPath)
		if svcErr != nil {
			http.Error(w, svcErr.Error(), http.StatusBadRequest)
		}
		// check whether a file exists at the given path
		_, err := os.Stat(res.FilePath)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, res.IndexPath)
			return
		} else if err != nil {
			logrus.Errorln(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// default to using service.FileServer to serve the static dir
		http.FileServer(http.Dir(res.StaticPath)).ServeHTTP(w, r)
	}
}

func (h Handler) ServiceHandler() http.HandlerFunc {
	serverProxy := h.Service.Server()
	return func(rw http.ResponseWriter, req *http.Request) {
		logrus.Infoln(req.URL.Path)
		serverProxy.Director(req)
		serverProxy.ServeHTTP(rw, req)
	}
}
