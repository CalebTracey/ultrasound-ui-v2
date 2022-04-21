package routes

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
	"ultrasound-client/go-server/log"
	"ultrasound-client/go-server/proxy"
)

const ServerUrlProd = "https://ultrasound-api.herokuapp.com"

type Handler struct {
	Service    proxy.FacadeI
	IndexPath  string
	StaticPath string
}

func (h Handler) InitializeRoutes() *mux.Router {
	//h.App = spaHandler{staticPath: "./web", indexPath: "index.html"}
	r := mux.NewRouter().StrictSlash(true)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(log.LoggingMiddleware)

	r.PathPrefix("/").Handler(h.ClientHandler())
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", h.ClientHandler()))
	r.PathPrefix("/api/").Handler(h.ServiceHandler())

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
	})
	return r
}

func (h Handler) ClientHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		staticIndexPath := filepath.Join(h.StaticPath, h.IndexPath)

		path, err := filepath.Abs(req.URL.Path)
		if err != nil {
			logrus.Errorln(err.Error())
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		servePath, svcErr := h.Service.Client(req.Context(), path)

		if svcErr != nil {
			http.Error(rw, svcErr.Error(), http.StatusInternalServerError)
			return
		}
		if strings.EqualFold(servePath, staticIndexPath) {
			http.ServeFile(rw, req, staticIndexPath)
		}

		http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(rw, req)

	}
}

func (h Handler) ServiceHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		targetUrl, err := url.Parse(ServerUrlProd)
		if err != nil {
			logrus.Errorln(err)
		}
		req.Host = targetUrl.Host
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		logrus.Infof("Service Proxy URL Path: %v", req.URL.Path)

		httputil.NewSingleHostReverseProxy(targetUrl).ServeHTTP(rw, req)
	}
}

func proxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "application/json")
		if req.Method == http.MethodOptions {
			return
		}
	}
}
