package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/ultra207/ultrasound-client/go-server/config"
	"gitlab.com/ultra207/ultrasound-client/go-server/internal/facade"
	"gitlab.com/ultra207/ultrasound-client/go-server/internal/routes"
	"gitlab.com/ultra207/ultrasound-client/go-server/internal/service"
	"os"
)

var configPath = "config.json"

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {
	defer deathScream()

	appConfig := config.NewConfigFromFile(configPath)
	appService, svcErr := facade.NewService(appConfig)
	if svcErr != nil {
		logrus.Panic(svcErr)
		panic(svcErr)
	}

	handler := routes.Handler{
		Service: appService,
	}

	env, envErr := appService.Environment()
	if envErr != nil {
		logrus.Errorf("environment error: %v", envErr.Error())
	}

	router := handler.InitializeRoutes()

	logrus.Infof("Current environment: %v", os.Getenv("ENVIRONMENT"))
	logrus.Fatal(service.ListenAndServe(env.Port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func deathScream() {
	if r := recover(); r != nil {
		logrus.Errorf("I panicked and am quitting: %v,", r)
		if err, ok := r.(stackTracer); ok {
			logrus.Tracef("%+v", err.StackTrace()[0:2])
		}
	}
}
