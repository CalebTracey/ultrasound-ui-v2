package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"os"
	"ultrasound-client/go-server/proxy"
)

const DefaultPort = "6088"
const StaticPath = "/web"
const IndexPath = "index.html"

func main() {
	defer deathScream()

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	logrus.Infof("Port is: %v", port)

	handler := Handler{
		Service:    proxy.NewService(),
		staticPath: StaticPath,
		indexPath:  IndexPath,
	}
	router := handler.initializeRoutes()
	logrus.Fatal(listenAndServe(port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func deathScream() {
	if r := recover(); r != nil {
		logrus.Errorf("I panicked and am quitting: %v,", r)
	}
}
