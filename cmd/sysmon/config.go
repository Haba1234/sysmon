package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

// Config структура конфигурации.
type Config struct {
	Logger     LoggerConf     // Logger - конфигурация логгера.
	Collection CollectionConf // Collection - Коллекция метрик.
}

// LoggerConf структура конфигурации.
type LoggerConf struct {
	Level string `config:"level"` // Level - уровень логирования
	Path  string `config:"path"`  // Path - путь к файлу лога.
}

// CollectionConf структура конфигурации.
type CollectionConf struct {
	LoadAverageEnabled bool `config:"loadAverage"` // LoadAverageEnabled - включение метрики 'load average'.
	BufSize            int  `config:"bufSize"`     // BufSize - глубина истории собираемых метрик.
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
