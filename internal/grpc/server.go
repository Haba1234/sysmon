package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	stat "github.com/Haba1234/sysmon/internal/grpc/api"
	"github.com/Haba1234/sysmon/internal/logger"
	"github.com/Haba1234/sysmon/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc -I ./api Service.proto --go_out=. --go-grpc_out=.

type Server struct {
	stat.UnimplementedStatisticsServer
	srv     *grpc.Server
	logg    *logger.Logger
	serv    *service.Collector
	bufSize int
}

func NewServer(logg *logger.Logger, serv *service.Collector, bufSize int) *Server {
	return &Server{
		logg:    logg,
		bufSize: bufSize,
		serv:    serv,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logg.Info("gRPC server " + addr + " starting...")
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer(grpc.StreamInterceptor(loggingServerInterceptor(s.logg)))
	stat.RegisterStatisticsServer(s.srv, s)

	if err := s.srv.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func loggingServerInterceptor(logger *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Info(fmt.Sprintf("method: %s, duration: %s, request: %+v", info.FullMethod, time.Since(time.Now()), srv))
		return handler(srv, ss)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("gRPC server stopping...")
	s.srv.GracefulStop()
	return nil
}

func (s *Server) ListStatistics(req *stat.SubscriptionRequest, stream stat.Statistics_ListStatisticsServer) error {
	n := req.GetPeriod().AsDuration()
	m := int(req.GetDepth())
	buf := time.Duration(s.bufSize) * time.Second
	s.logg.Debug(fmt.Sprintf("Client connected - period: %v, depth: %vs, buf: %v", n, m, buf), "grpc")

	if n <= 0 {
		return status.Error(
			codes.InvalidArgument,
			"Period must be greater than 0 sec",
		)
	}
	if n > buf {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Period must be less than %v sec", buf),
		)
	}
	if m <= 0 {
		return status.Error(
			codes.InvalidArgument,
			"Depth must be greater than 0 sec",
		)
	}
	if m > s.bufSize {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Depth must be less than %v sec", s.bufSize),
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
			if s.serv.StatsData.Counter >= m { // Данные накопились, можно отправлять.
				if err := stream.Send(s.getStatistics(m)); err != nil {
					return err
				}
			}
		}
	}
}

func (s *Server) getStatistics(m int) *stat.StatisticsResponse {
	stats := &stat.StatisticsResponse{
		Status: "OK",
	}
	s.serv.ReadStats(m)

	data := s.serv.StatsData
	stats.La = &stat.LoadAverage{
		OneMin:     data.La[0],
		FiveMin:    data.La[1],
		FifteenMin: data.La[2],
	}
	stats.Cp = &stat.CPUAverage{
		UserMode: data.CPU[0],
		SysMode:  data.CPU[1],
		Idle:     data.CPU[2],
	}
	return stats
}
