package facade

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"ultrasound-client/go-server/config"
)

type Facade interface {
	Server() *httputil.ReverseProxy
}

type Service struct {
	clientUrl string
	serverUrl string
}

func NewService(appConfig *config.Config) Service {
	return Service{
		clientUrl: appConfig.ClientUrl,
		serverUrl: appConfig.ServerUrl,
	}
}

func (s Service) Server() *httputil.ReverseProxy {
	target, svrErr := url.Parse(s.serverUrl)
	if svrErr != nil {
		panic(svrErr)
	}
	origin, uiErr := url.Parse(s.clientUrl)
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
