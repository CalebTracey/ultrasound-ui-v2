package server

import "ultrasound-client/go-server/config"

type Service struct {
	ServerUrl string
	Port      string
}

func InitializeServerService(appConfig *config.Config) *Service {
	clientConf := appConfig.ServerConfig
	return &Service{
		ServerUrl: clientConf.Url,
		Port:      appConfig.Port,
	}
}
