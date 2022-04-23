package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ultrasound-client/go-server/facade"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

type Handler struct {
	Service facade.Facade
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(h.loggingMiddleware)
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
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// prepend the path with the path to the static directory
		res := h.Service.Client(path)

		// check whether a file exists at the given path
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, res.FilePath)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// default to using service.FileServer to serve the static dir
		http.FileServer(http.Dir(res.IndexPath)).ServeHTTP(w, r)
	}
}

func (h Handler) ServiceHandler() http.HandlerFunc {
	serverProxy := proxyRequestHandler(h.Service.Server())
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

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}
