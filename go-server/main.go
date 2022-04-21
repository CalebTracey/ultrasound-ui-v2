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
	"ultrasound-client/go-server/proxy"
)

const DefaultPort = ":6088"

func main() {
	defer deathScream()

	port := os.Getenv("PORT")

	if port == "" {
		logrus.Fatal("$PORT must be set")
	}
	handler := Handler{
		Service: proxy.NewService(),
	}
	router := handler.initializeRoutes()
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
