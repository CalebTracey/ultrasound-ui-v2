package config

type Config struct {
	Env          string
	Port         string
	ClientConfig ClientConfig
	ServerConfig ServerConfig
	Hash         string
}

type ClientConfig struct {
	Name string
	Url  string
}

type ServerConfig struct {
	Name string
	Url  string
}
