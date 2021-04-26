package sysmon

import "context"

// gRPC Server.
type GRPCServer interface {
	Start(ctx context.Context, addr string) error
	Stop(ctx context.Context)
}

// Logger.
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg, pkg string)
}
