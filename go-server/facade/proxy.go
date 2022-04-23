package facade

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ultrasound-client/go-server/config"
	"gitlab.com/ultra207/ultrasound-client/go-server/service/client"
	"gitlab.com/ultra207/ultrasound-client/go-server/service/server"

	"gitlab.com/ultra207/ultrasound-client/go-server/service/environment"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

type ProxyFacade interface {
	Client(urlPath string) (ClientResponse, error)
	Server() *httputil.ReverseProxy
	Environment() environment.Response
}

type Service struct {
	ClientSvc client.ServiceI
	ServerSvc server.ServiceI
	EnvSvc    environment.ServiceI
}

func NewService(appConfig *config.Config) Service {
	clientSvc := client.InitializeClientSvc(appConfig)
	serverSvc := server.InitializeServerSvc(appConfig)
	envSvc := environment.InitializeEnvService(appConfig)
	return Service{
		ClientSvc: &clientSvc,
		ServerSvc: &serverSvc,
		EnvSvc:    &envSvc,
	}
}

func (s Service) Client(urlPath string) (ClientResponse, error) {
	svc := s.ClientSvc.Client()

	path, err := filepath.Abs(urlPath)
	if err != nil {
		logrus.Errorln(err.Error())
		return ClientResponse{}, err
	}
	return ClientResponse{
		// prepend the path with the path to the static directory
		FilePath: filepath.Join(svc.StaticPath, path),
		// path to index.html for when file does not exist
		IndexPath:  filepath.Join(svc.StaticPath, svc.IndexPath),
		StaticPath: svc.StaticPath,
	}, nil
}

func (s Service) Server() *httputil.ReverseProxy {
	svc := s.ServerSvc.Server()

	target, svrErr := url.Parse(svc.TargetUrl)
	if svrErr != nil {
		panic(svrErr)
	}
	origin, uiErr := url.Parse(svc.HostUrl)
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

func (s Service) Environment() environment.Response {
	envResponse, envErr := s.EnvSvc.Set()
	if envErr != nil {
		logrus.Panic(envErr.Error())
	}
	return envResponse
}
