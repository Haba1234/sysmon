package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

// Config структура конфигурации:
// Logger - конфигурация логгера;
// Collection - Коллекция метрик.
type Config struct {
	Logger     LoggerConf
	Collection CollectionConf
}

// LoggerConf содержит:
// Level - уровень логирования
// Path - путь к файлу лога.
type LoggerConf struct {
	Level string `config:"level"`
	Path  string `config:"path"`
}

// CollectionConf содержит:
// LoadAverageEnabled - включение метрики 'load average'
// BufSize - глубина истории собираемых метрик.
type CollectionConf struct {
	LoadAverageEnabled bool `config:"loadAverage"`
	BufSize            int  `config:"bufSize"`
}

func NewConfig(path string) (*Config, error) {
	// default values
	cfg := Config{
		Logger: LoggerConf{
			Level: "INFO",
			Path:  "/tmp/logfile.log",
		},
		Collection: CollectionConf{
			LoadAverageEnabled: true,
			BufSize:            60,
		},
	}

	loader := confita.NewLoader(
		file.NewBackend(path),
	)
	if err := loader.Load(context.Background(), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
