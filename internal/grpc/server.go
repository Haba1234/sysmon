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
	service *service.LoadAverage
	logg    *logger.Logger
}

func NewServer(logg *logger.Logger, service *service.LoadAverage) *Server {
	return &Server{
		service: service,
		logg:    logg,
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
	s.logg.Debug(fmt.Sprintf("Client connected - period: %v, depth: %v", req.GetPeriod().AsDuration(), req.GetDepth()), "grpc")
	ticker := time.NewTicker(req.GetPeriod().AsDuration())
	s.logg.Warn("Старт таймера с заданной периодичностью")
loop:
	for {
		select {
		case <-stream.Context().Done():
			s.logg.Info("gRPC client stopped")
			ticker.Stop()
			break loop
		case <-ticker.C:
			result, err := s.service.AverageValue(int(req.GetDepth()))
			if err != nil {
				s.logg.Error(fmt.Sprintf("load average request error: %v", err))
				return status.Errorf(codes.Internal, "load average request error: %v", err)
			}
			loadAvg := ""
			for _, s := range result {
				loadAvg = loadAvg + " " + s
			}
			s.logg.Info(fmt.Sprintf("Средняя заргузка: %s", loadAvg))
			la := stat.LoadAverage{
				OneMin:     result[0],
				FiveMin:    result[1],
				FifteenMin: result[2],
			}

			if err := stream.Send(&stat.StatisticsResponse{La: &la}); err != nil {
				return err
			}
		}
	}
	return nil
}
