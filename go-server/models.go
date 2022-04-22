package main

import (
	"bytes"
	"net/http"
	"ultrasound-client/go-server/proxy"
)

type Handler struct {
	Service    proxy.Facade
	staticPath string
	indexPath  string
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}
