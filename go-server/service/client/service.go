package client

import "ultrasound-client/go-server/config"

type Service struct {
	ClientUrl string
	Port      string
}

func InitializeClientService(appConfig *config.Config) *Service {
	clientConf := appConfig.ClientConfig
	return &Service{
		ClientUrl: clientConf.Url,
		Port:      appConfig.Port,
	}
}
