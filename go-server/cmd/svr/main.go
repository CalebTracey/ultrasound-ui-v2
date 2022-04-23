package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ultrasound-client/go-server/config"
	"gitlab.com/ultra207/ultrasound-client/go-server/facade"
	"gitlab.com/ultra207/ultrasound-client/go-server/routes"
	"gitlab.com/ultra207/ultrasound-client/go-server/service"
	"os"
)

var configPath = "config.json"

func main() {
	defer deathScream()

	appConfig := config.NewConfigFromFile(configPath)

	handler := routes.Handler{
		Service: facade.NewService(appConfig),
	}

	env := handler.Service.Environment()
	router := handler.InitializeRoutes()

	logrus.Infof("Current environment: %v", os.Getenv("ENVIRONMENT"))
	logrus.Fatal(service.ListenAndServe(env.Port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func deathScream() {
	if r := recover(); r != nil {
		logrus.Errorf("I panicked and am quitting: %v,", r)
	}
}
