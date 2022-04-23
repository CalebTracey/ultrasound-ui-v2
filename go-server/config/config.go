package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func New(configPath string) (*Config, []error) {
	return newConfig(&defaultBuilder{}, configPath)
}

func NewConfigFromFile(configPath string) *Config {
	logrus.Infoln(configPath)
	conf, confErrs := New(configPath)
	if len(confErrs) > 0 || conf == nil {
		for _, err := range confErrs {
			panic(fmt.Sprintf("Unable to load config: %v\n", err.Error()))
		}
		if conf == nil {
			panic("Config File returned nil")
		}
		panic("Exiting: Could not load config file")
	}
	return conf
}

func newConfig(builder configBuilder, configPath string) (*Config, []error) {
	var err error

	configFile, err := builder.Load(configPath)
	if err != nil {
		return nil, []error{err}
	}
	defer func(configFile *os.File) {
		closeErr := configFile.Close()
		if closeErr != nil {
			logrus.Errorln(closeErr.Error())
		}
	}(configFile)

	err = builder.Read(configFile)
	if err != nil {
		return nil, []error{}
	}

	return builder.GetConfig(), nil
}
