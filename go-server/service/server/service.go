package server

import "gitlab.com/ultra207/ultrasound-client/go-server/config"

type ServiceI interface {
	Server() Response
}
type Service struct {
	hostUrl   string
	targetUrl string
}

type Response struct {
	HostUrl   string
	TargetUrl string
}

func InitializeServerSvc(appConfig *config.Config) Service {
	return Service{
		hostUrl:   appConfig.ClientConfig.Url,
		targetUrl: appConfig.ClientConfig.Url,
	}
}

func (s *Service) Server() Response {
	return Response{
		HostUrl:   s.hostUrl,
		TargetUrl: s.targetUrl,
	}
}
