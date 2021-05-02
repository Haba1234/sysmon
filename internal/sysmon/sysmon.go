package sysmon

import "context"

// gRPC Server интерфейсы.
type GRPCServer interface {
	Start(ctx context.Context, addr string) error
	Stop(ctx context.Context)
}

// Logger интерфейсы.
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg, pkg string)
}
