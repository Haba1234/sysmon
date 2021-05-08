package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/Haba1234/sysmon/internal/grpc"
	"github.com/Haba1234/sysmon/internal/logger"
	"github.com/Haba1234/sysmon/internal/service"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	var portServer string
	flag.StringVar(&portServer, "port", "8080", "gRPC server port number")
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	settings := service.Collection{
		LoadAverageEnabled: config.Collection.LoadAverageEnabled,
		CPUEnabled:         config.Collection.CPUEnabled,
		BufSize:            config.Collection.BufSize,
	}
	collector := service.NewCollector(logg, settings)

	server := grpc.NewServer(logg, collector, config.Collection.BufSize)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop gRPC server: " + err.Error())
		}
	}()

	logg.Info("system monitoring is running...")

	go func() {
		if err := collector.Start(ctx); err != nil {
			logg.Error("failed to start 'service collector' service: " + err.Error())
			cancel()
		}

		addrServer := net.JoinHostPort("", portServer)
		if err := server.Start(ctx, addrServer); err != nil {
			logg.Error("failed to start gRPC server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()
	logg.Warn("system monitoring stopped.")
}
