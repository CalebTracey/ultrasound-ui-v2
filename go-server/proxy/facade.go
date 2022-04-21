package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

//const clientURL = "http://localhost:6088"
const clientURLProd = "https://ultrasound-ui.herokuapp.com"

//const serverURL = "http://localhost:6080"

type ProxyFacade interface {
	//Client() *httputil.ReverseProxy
	Server() *httputil.ReverseProxy
}

func NewService() Service {
	return Service{}
}

type Service struct{}

func (s Service) Server() *httputil.ReverseProxy {
	origin, _ := url.Parse(clientURLProd)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}
	proxy := &httputil.ReverseProxy{Director: director}

	return proxy
}
