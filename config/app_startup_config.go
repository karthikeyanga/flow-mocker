package config

import (
	"mocker/common"
	"mocker/util/log"
)

type AppStartupConfigType struct {
	AppServerConfig AppServerConfig    `yaml:"AppServerConfig"`
	DatabaseConfig  DatabaseConfigType `yaml:"DatabaseConfig"`
	LoggingConfig   LoggingConfigType  `yaml:"LoggingConfig"`
	ServerConfig    ServerConfigType   `yaml:"ServerConfig"`
	APIMode         common.APIModeType `yaml:"APIMode"`
	AppName         string             `yaml:"AppName"`
	Constants       map[string]string  `yaml:"Constants"`
}
type DatabaseConfigType struct {
	Dialect                string
	Path                   string
	Host                   string
	Port                   string
	UserName               string
	Password               string
	Database               string
	Protocol               string
	Timeout                int
	MaxIdleConnectionCount int
	MaxConnectionCount     int
	LogFile                string
}

type LoggingConfigType struct {
	Level log.LogLevel `yaml:"LogLevel"`
	Path  string       `yaml:"Path"`
}

type ServerConfigType struct {
	Host       string `yaml:"Host"`
	Port       string `yaml:"Port"`
	AccessLog  string `yaml:"AccessLog"`
	ErrorLog   string `yaml:"ErrorLog"`
	StaticPath string `yaml:"StaticPath"`
}

type AppServerConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
}
