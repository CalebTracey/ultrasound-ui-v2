package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"ultrasound-client/go-server/proxy"
)

const infoLevel = logrus.InfoLevel

type Handler struct {
	Service proxy.ProxyFacade
	App     spaHandler
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h Handler) initializeRoutes() *mux.Router {
	h.App = spaHandler{staticPath: "./web", indexPath: "index.html"}
	fs := http.FileServer(http.Dir("./web"))

	r := mux.NewRouter().StrictSlash(true)
	r.Use(h.loggingMiddleware)

	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", fs))
	r.PathPrefix("/").Handler(fs)

	r.Handle("/api/", h.ServiceHandler())
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			return
		}
	})
	return r
}

func (h Handler) ClientHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		logrus.Infoln("Service Handler")
		logrus.Infoln(req.URL.Path)
		http.FileServer(http.Dir(h.App.staticPath)).ServeHTTP(rw, req)
	}
}

func (h Handler) ServiceHandler() http.HandlerFunc {
	serverProxy := proxyRequestHandler(h.Service.Server())
	logrus.Infoln("Service Handler")
	return serverProxy
}

func proxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		logrus.Infoln(req.URL.Path)
		proxy.Director(req)
		proxy.ServeHTTP(rw, req)
	}
}

func (h Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logrus.Errorf("Error reading body: %v", err)
			http.Error(rw, "can't read body", http.StatusBadRequest)
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		mrw := &MyResponseWriter{
			ResponseWriter: rw,
			buf:            &bytes.Buffer{},
		}

		next.ServeHTTP(mrw, req)

		if _, ioErr := io.Copy(rw, mrw.buf); err != nil {
			logrus.Errorf("Failed to send out response: %v", ioErr)
		}
	})
}
