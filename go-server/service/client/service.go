package client

import (
	"gitlab.com/ultra207/ultrasound-client/go-server/config"
)

type ServiceI interface {
	Client() Response
}
type Service struct {
	staticPath string
	indexPath  string
}

type Response struct {
	StaticPath string
	IndexPath  string
}

func InitializeClientSvc(appConfig *config.Config) Service {
	return Service{
		staticPath: appConfig.ClientConfig.StaticPath,
		indexPath:  appConfig.ClientConfig.IndexPath,
	}
}

func (s *Service) Client() Response {
	return Response{
		StaticPath: s.staticPath,
		IndexPath:  s.indexPath,
	}
}
