package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}
