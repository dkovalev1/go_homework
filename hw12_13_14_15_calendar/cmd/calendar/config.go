package main

import (
	"log"

	"github.com/BurntSushi/toml" //nolint
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	REST    RESTConf
	GRPC    GRPCConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type    string
	Connstr string
}

type RESTConf struct {
	Port int
}

type GRPCConf struct {
	Port int
}

func NewConfig(configFile string) Config {
	var config Config
	// Read the TOML file

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	log.Printf("log level: %s", config.Logger.Level)

	return config
}

// TODO
