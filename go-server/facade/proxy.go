package facade

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ultrasound-client/go-server/config"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

type ProxyFacade interface {
	Client(urlPath string) (ClientResponse, error)
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

func (s Service) Client(urlPath string) (ClientResponse, error) {
	path, err := filepath.Abs(urlPath)
	if err != nil {
		logrus.Errorln(err.Error())
		return ClientResponse{}, err
	}
	return ClientResponse{
		// prepend the path with the path to the static directory
		FilePath: filepath.Join(s.staticPath, path),
		// path to index.html for when file does not exist
		IndexPath:  filepath.Join(s.staticPath, s.indexPath),
		StaticPath: s.staticPath,
	}, nil
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
