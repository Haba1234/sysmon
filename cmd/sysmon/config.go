package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Logger     LoggerConf
	Collection CollectionConf
}

type LoggerConf struct {
	Level string `config:"level"`
	Path  string `config:"path"`
}

type CollectionConf struct {
	LoadAverageEnabled bool  `config:"loadAverage"`
	BufSize            int64 `config:"bufSize"`
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
