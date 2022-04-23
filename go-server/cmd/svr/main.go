package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"os"
	"ultrasound-client/go-server/config"
	"ultrasound-client/go-server/facade"
	"ultrasound-client/go-server/routes"
	"ultrasound-client/go-server/service"
)

const DefaultPort = "6088"
const StaticPath = "/web"
const IndexPath = "index.html"
const ServerUrlProd = "https://ultrasound-api.herokuapp.com"
const ClientUrlProd = "https://ultrasound-ui.herokuapp.com"

func main() {
	defer deathScream()

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	logrus.Infof("Port is: %v", port)

	appConfig := &config.Config{
		ClientUrl: ClientUrlProd,
		ServerUrl: ServerUrlProd,
	}

	handler := routes.Handler{
		Service:    facade.NewService(appConfig),
		StaticPath: StaticPath,
		IndexPath:  IndexPath,
	}
	router := handler.InitializeRoutes()

	logrus.Fatal(service.ListenAndServe(port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func deathScream() {
	if r := recover(); r != nil {
		logrus.Errorf("I panicked and am quitting: %v,", r)
	}
}
