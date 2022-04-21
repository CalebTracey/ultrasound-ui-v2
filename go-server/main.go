package main

import (
	"context"
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

//
//type Handler struct {
//	staticPath string
//	indexPath  string
//}

const Port = "0.0.0.0:6088"
const StaticPath = "build"
const IndexPath = "index.html"

func main() {
	defer deathScream()

	handler := Handler{
		Service: proxy.NewService(),
	}

	router := handler.initializeRoutes()
	logrus.Fatal(ListenAndServe(Port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func ListenAndServe(addr string, handler http.Handler) error {
	logrus.Infof("=== HTTP Server address %v ===", addr)

	srv := &http.Server{
		Addr:         addr,
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
