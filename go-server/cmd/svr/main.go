package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ult-config/config"
	"os"
	"ultrasound-client/go-server/facade"
	"ultrasound-client/go-server/routes"
	"ultrasound-client/go-server/service"
)

const DefaultPort = "6088"

var configPath = "config.json"

func main() {
	defer deathScream()

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	logrus.Infof("Port is: %v", port)

	appConfig := config.NewConfigFromFile(configPath)

	handler := routes.Handler{
		Service: facade.NewService(appConfig),
	}

	router := handler.InitializeRoutes()
	logrus.Fatal(service.ListenAndServe(port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func deathScream() {
	if r := recover(); r != nil {
		logrus.Errorf("I panicked and am quitting: %v,", r)
	}
}
