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

const DefaultPort = "6088"

var configPath = "config.json"

func main() {
	defer deathScream()

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	appConfig := config.NewConfigFromFile(configPath)
	// set ENVIRONMENT env variable for react to set api url in src/service/api.tsx
	envErr := os.Setenv("ENVIRONMENT", appConfig.Env)
	if envErr != nil {
		logrus.Errorf("failed to set ENVIRONMENT; err: %v", envErr.Error())
	}

	logrus.Infof("Current environment: %v", os.Getenv("ENVIRONMENT"))

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
