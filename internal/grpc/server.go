package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	stat "github.com/Haba1234/sysmon/internal/grpc/api"
	"github.com/Haba1234/sysmon/internal/sysmon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc -I ./api Service.proto --go_out=. --go-grpc_out=.

type Server struct {
	mu *sync.Mutex
	stat.UnimplementedStatisticsServer
	srv       *grpc.Server
	logg      sysmon.Logger
	collector sysmon.Collector
	bufSize   int
	maxPeriod time.Duration
}

const twoMin = 120

func NewServer(logg sysmon.Logger, collector sysmon.Collector, bufSize int) *Server {
	return &Server{
		mu:        &sync.Mutex{},
		logg:      logg,
		bufSize:   bufSize,
		collector: collector,
		maxPeriod: time.Duration(twoMin) * time.Second,
	}
}

// Start запуск сервера gRPC.
func (s *Server) Start(ctx context.Context, addr string) error {
	s.logg.Info("gRPC server " + addr + " running...")
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.srv = grpc.NewServer(grpc.StreamInterceptor(loggingServerInterceptor(s.logg)))
	s.mu.Unlock()
	stat.RegisterStatisticsServer(s.srv, s)

	if err := s.srv.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func loggingServerInterceptor(logg sysmon.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logg.Info(fmt.Sprintf("method: %s, duration: %s, request: %+v", info.FullMethod, time.Since(time.Now()), srv))
		return handler(srv, ss)
	}
}

// Stop останов сервера gRPC.
func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("gRPC server stopped")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.srv.GracefulStop()
	return nil
}

func (s *Server) ListStatistics(req *stat.SubscriptionRequest, stream stat.Statistics_ListStatisticsServer) error {
	n := req.GetPeriod().AsDuration()
	m := int(req.GetDepth())
	buf := time.Duration(s.bufSize) * time.Second
	s.logg.Debug(fmt.Sprintf("client connected - period: %v, depth: %vs, buf: %v", n, m, buf), "grpc")

	if n <= 0 {
		return status.Error(
			codes.InvalidArgument,
			"period must be greater than 0 sec",
		)
	}
	if n > s.maxPeriod {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("period must be less than %v sec", s.maxPeriod),
		)
	}
	if m <= 0 {
		return status.Error(
			codes.InvalidArgument,
			"depth must be greater than 0 sec",
		)
	}
	if m > s.bufSize {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("depth must be less than %v sec", s.bufSize),
		)
	}

	ticker := time.NewTicker(n)
	for {
		select {
		case <-stream.Context().Done():
			s.logg.Info("gRPC client stopped")
			ticker.Stop()
			return nil
		case <-ticker.C:
			statusServices := s.collector.GetStatusServices()
			if statusServices.La.StatusCode != sysmon.ServiceRun && statusServices.CPU.StatusCode != sysmon.ServiceRun {
				// Нет ни одного запущенного сервиса!
				return status.Error(
					codes.Aborted,
					"all services are stopped",
				)
			}
			if statusServices.La.Counter < m && statusServices.CPU.Counter < m {
				// Данные еще не готовы для отправки
				break
			}
			if err := stream.Send(s.getStatistics(m)); err != nil {
				return err
			}
		}
	}
}

func (s *Server) getStatistics(m int) *stat.StatisticsResponse {
	stats := &stat.StatisticsResponse{
		Status: "OK",
	}
	result := s.collector.GetStats(m)

	stats.La = &stat.LoadAverage{
		Status:     statusCodeToStr(result.La.StatusCode),
		OneMin:     result.La.Data[0],
		FiveMin:    result.La.Data[1],
		FifteenMin: result.La.Data[2],
	}
	stats.Cp = &stat.CPUAverage{
		Status: statusCodeToStr(result.CPU.StatusCode),
		User:   result.CPU.Data[0],
		Sys:    result.CPU.Data[1],
		Idle:   result.CPU.Data[2],
	}
	return stats
}

func statusCodeToStr(code int) string {
	switch code {
	case 1:
		return "works"
	case 2:
		return "stopped"
	case 3:
		return "stopped with errors"
	default:
		return "unknown"
	}
}
