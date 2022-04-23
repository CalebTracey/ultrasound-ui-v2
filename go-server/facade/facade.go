package facade

import (
	"gitlab.com/ultra207/ultrasound-client/go-server/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

type Facade interface {
	Client(path string) ClientResponse
	Server() *httputil.ReverseProxy
}

type Service struct {
	staticPath string
	indexPath  string
	clientUrl  string
	serverUrl  string
}

func NewService(appConfig *config.Config) Service {
	return Service{
		staticPath: appConfig.ClientConfig.StaticPath,
		indexPath:  appConfig.ClientConfig.IndexPath,
		clientUrl:  appConfig.ClientConfig.Url,
		serverUrl:  appConfig.ServerConfig.Url,
	}
}

func (s Service) Client(path string) ClientResponse {
	return ClientResponse{
		FilePath:  filepath.Join(s.staticPath, path),
		IndexPath: filepath.Join(s.staticPath, s.indexPath),
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
