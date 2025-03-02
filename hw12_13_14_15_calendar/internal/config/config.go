package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml" //nolint
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Storage  StorageConf
	RabbitMQ SchedulerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

type StorageConf struct {
	Type    string
	Connstr string
	// TODO
}

type SchedulerConf struct {
	Connstr    string
	Interval   time.Duration
	KeepEvents time.Duration
}

func NewConfig(configFile string) *Config {
	var config Config
	// Read the TOML file

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	log.Printf("log level: %s", config.Logger.Level)

	return &config
}
