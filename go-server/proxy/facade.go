package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const serverUrlProd = "https://ultrasound-api.herokuapp.com"
const clientUrlProd = "https://ultrasound-ui.herokuapp.com"

type Facade interface {
	Server() *httputil.ReverseProxy
}

func NewService() Service {
	return Service{}
}

type Service struct{}

func (s Service) Server() *httputil.ReverseProxy {
	target, svrErr := url.Parse(serverUrlProd)
	if svrErr != nil {
		panic(svrErr)
	}
	origin, uiErr := url.Parse(clientUrlProd)
	if uiErr != nil {
		panic(uiErr)
	}
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", target.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path

	}
	proxy := &httputil.ReverseProxy{Director: director}

	return proxy
}
