package config

import (
	bytes2 "bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

type configBuilder interface {
	Load(string) (*os.File, error)
	Read(io.Reader) error
	GetConfig() *Config
	GetConfigPath() string
}

type defaultBuilder struct {
	config     *Config
	configPath string
}

func (b *defaultBuilder) GetConfig() *Config {
	return b.config
}

func (b *defaultBuilder) GetConfigPath() string {
	return b.configPath
}

func (b *defaultBuilder) Load(path string) (*os.File, error) {
	logrus.Tracef("Loading config: %v", path)
	b.configPath = path

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file %v; %v", path, err.Error())
	}

	return file, err
}

func (b *defaultBuilder) Read(configData io.Reader) error {
	logrus.Trace("Reading config data")

	config, errs := buildInitialConfig(configData)
	if errs != nil {
		return errs
	}

	b.config = config
	return nil
}

func buildInitialConfig(configData io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(configData)
	if err != nil {
		return nil, fmt.Errorf("error reading config data: %v", err.Error())
	}

	br := bytes2.NewReader(bytes)

	c := &Config{}
	decoder := json.NewDecoder(br)
	decoder.DisallowUnknownFields()
	decodeErr := decoder.Decode(&c)
	if decodeErr != nil {
		return nil, fmt.Errorf("error decoding config data: %v", decodeErr)
	}

	c.Hash = fmt.Sprintf("%x", md5.Sum(bytes))

	return c, nil
}
