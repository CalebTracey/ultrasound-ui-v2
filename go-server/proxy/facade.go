package proxy

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"ultrasound-client/go-server/config"
	"ultrasound-client/go-server/service/client"
	"ultrasound-client/go-server/service/server"
)

const StaticPath = "./web"
const IndexPath = "index.html"

type FacadeI interface {
	Client(ctx context.Context, path string) (string, error)
}

type Service struct {
	ClientSvc  *client.Service
	ServerSvc  *server.Service
	StaticPath string
	IndexPath  string
}

func NewService(appConfig *config.Config) Service {
	return Service{
		ClientSvc:  client.InitializeClientService(appConfig),
		ServerSvc:  server.InitializeServerService(appConfig),
		StaticPath: StaticPath,
		IndexPath:  IndexPath,
	}
}

func (s Service) Client(ctx context.Context, path string) (string, error) {
	path = filepath.Join(s.StaticPath, path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		staticIndexPath := filepath.Join(s.StaticPath, s.IndexPath)
		return staticIndexPath, nil
	} else if err != nil {
		logrus.Errorln(err.Error())
		return "", err
	}
	return "", nil
}
