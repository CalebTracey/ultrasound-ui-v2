package log

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func LoggingMiddleware(next http.Handler) http.Handler {
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
