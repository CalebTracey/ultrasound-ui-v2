package main

import (
	"context"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"ultrasound-client/go-server/config"
	"ultrasound-client/go-server/proxy"
	"ultrasound-client/go-server/routes"
)

const DefaultPort = "6088"

var (
	configPath = "local_config.json"
)

func main() {
	defer deathScream()

	appConfig := config.NewConfigFromFile(configPath)
	proxyService := proxy.NewService(appConfig)
	port := os.Getenv("PORT")

	if port == "" {
		port = DefaultPort
		logrus.Infof("PORT is: %v", port)
	}

	handler := routes.Handler{
		Service:    &proxyService,
		StaticPath: proxyService.StaticPath,
		IndexPath:  proxyService.IndexPath,
	}

	router := handler.InitializeRoutes()
	logrus.Fatal(listenAndServe(port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func listenAndServe(addr string, handler http.Handler) error {
	logrus.Infof("=== HTTP Server address %v ===", addr)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", addr),
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	wgServer := sync.WaitGroup{}
	wgServer.Add(2)

	var serverError error

	go func() {
		defer wgServer.Done()
		killSignal := <-signals
		switch killSignal {
		case os.Interrupt:
			logrus.Infoln("SIGINT recieved (Control-C ?)")
		case syscall.SIGTERM:
			logrus.Infoln("SIGTERM recieved (Heroku shutdown?)")
		case nil:
			return
		}
		logrus.Infoln("graceful shutdown...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			logrus.Error(err.Error())
		}
		logrus.Infoln("graceful shutdown complete")
	}()

	go func() {
		defer wgServer.Done()
		if err := srv.ListenAndServe(); err != nil {
			serverError = err
		}
		signals <- nil
	}()
	wgServer.Wait()
	return serverError
}

func deathScream() {
	if r := recover(); r != nil {
		logrus.Errorf("I panicked and am quitting: %v,", r)
	}
}
