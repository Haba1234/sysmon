package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	grpc "github.com/Haba1234/sysmon/internal/grpc/api"
	"github.com/Haba1234/sysmon/internal/mocks"
	"github.com/Haba1234/sysmon/internal/sysmon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Suite struct {
	suite.Suite
	log       *mocks.Logger
	collector *mocks.Collector
	bufSize   int
	la        []string
	cpu       []string
	m         int
}

func (s *Suite) SetupSuite() {
	s.log = new(mocks.Logger)
	s.log.On("Info", mock.Anything).Maybe()
	s.log.On("Debug", mock.Anything, "grpc").Maybe()

	s.m = 1
	s.la = []string{"1.0", "2.0", "3.0"}
	s.cpu = []string{"1.0", "2.0", "97.0"}

	s.collector = new(mocks.Collector)
	s.collector.On("GetStats", s.m).Return(sysmon.StatisticsData{
		La: sysmon.OutputData{
			Data:       s.la,
			Counter:    1,
			StatusCode: sysmon.ServiceRun,
		},
		CPU: sysmon.OutputData{
			Data:       s.cpu,
			Counter:    1,
			StatusCode: sysmon.ServiceStop,
		},
	}).Maybe()

	s.bufSize = 10
}

func (s *Suite) TearDownTest() {
	s.log.AssertExpectations(s.T())
	s.collector.AssertExpectations(s.T())
}

func (s *Suite) TestStartStop() {
	defer goleak.VerifyNone(s.T())
	server := NewServer(s.log, s.collector, s.bufSize)

	go func() {
		err := server.Start(context.Background(), ":8080")
		s.NoError(err)
	}()

	time.Sleep(10 * time.Millisecond)
	err := server.Stop(context.Background())
	s.NoError(err)
}

func (s *Suite) TestGetStatistics() {
	server := NewServer(s.log, s.collector, s.bufSize)
	result := server.getStatistics(s.m)
	s.Equal("works", result.GetLa().GetStatus())
	s.Equal(s.la[0], result.GetLa().GetOneMin())
	s.Equal(s.la[1], result.GetLa().GetFiveMin())
	s.Equal(s.la[2], result.GetLa().GetFifteenMin())
	s.Equal("stopped", result.GetCp().GetStatus())
	s.Equal(s.cpu[0], result.GetCp().GetUser())
	s.Equal(s.cpu[1], result.GetCp().GetSys())
	s.Equal(s.cpu[2], result.GetCp().GetIdle())
}

func (s *Suite) TestListStatisticsConditions() {
	server := NewServer(s.log, s.collector, s.bufSize)

	tests := []struct {
		period *durationpb.Duration
		depth  int64
		err    string
	}{
		{
			period: durationpb.New(time.Duration(5) * time.Second),
			depth:  int64(s.bufSize + 1),
			err:    fmt.Sprintf("depth must be less than %v sec", s.bufSize),
		},
		{
			period: durationpb.New(time.Duration(5) * time.Second),
			depth:  0,
			err:    "depth must be greater than 0 sec",
		},
		{
			period: durationpb.New(server.maxPeriod + 1),
			depth:  5,
			err:    fmt.Sprintf("period must be less than %v sec", server.maxPeriod),
		},
		{
			period: durationpb.New(time.Duration(0) * time.Second),
			depth:  5,
			err:    "period must be greater than 0 sec",
		},
	}
	stream := new(mocks.Statistics_ListStatisticsServer)

	for _, tt := range tests {
		tt := tt
		s.Run("threshold values", func() {
			s.T().Parallel()
			req := &grpc.SubscriptionRequest{
				Period: tt.period,
				Depth:  tt.depth,
			}

			stream.On("Context").Return(context.WithTimeout(context.Background(), time.Millisecond))
			err := server.ListStatistics(req, stream)
			s.Contains(err.Error(), tt.err)
		})
	}
	stream.AssertExpectations(s.T())
}

func (s *Suite) TestListStatisticsOneTickDone() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	stream := new(mocks.Statistics_ListStatisticsServer)
	stream.On("Context").Return(context.WithTimeout(context.Background(), 1010*time.Millisecond))
	stream.On("Send", mock.Anything).Return(nil)

	stat := s.structureFillStatusServices(1, sysmon.ServiceRun, sysmon.ServiceRun)
	s.collector.On("GetStatusServices").Return(stat)

	server := NewServer(s.log, s.collector, s.bufSize)
	req := s.structureFillSubscriptionRequest(1, 1)
	err := server.ListStatistics(req, stream)
	s.NoError(err)

	stream.AssertExpectations(s.T())
}

func (s *Suite) TestListStatisticsOneTickAborted() {
	s.T().Skip()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	stream := new(mocks.Statistics_ListStatisticsServer)
	stream.On("Context").Return(context.WithTimeout(context.Background(), 1010*time.Millisecond))

	stat := s.structureFillStatusServices(1, sysmon.ServiceStop, sysmon.ServiceError)
	s.collector.On("GetStatusServices").Return(stat)

	server := NewServer(s.log, s.collector, s.bufSize)
	req := s.structureFillSubscriptionRequest(1, 1)
	err := server.ListStatistics(req, stream)
	s.Error(err)
	stream.AssertExpectations(s.T())
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) structureFillStatusServices(count, codeLA, codeCPU int) sysmon.StatusServices {
	return sysmon.StatusServices{
		La: sysmon.Status{
			Counter:    count,
			StatusCode: codeLA,
		},
		CPU: sysmon.Status{
			Counter:    count,
			StatusCode: codeCPU,
		},
	}
}

func (s *Suite) structureFillSubscriptionRequest(period, depth int) *grpc.SubscriptionRequest {
	return &grpc.SubscriptionRequest{
		Period: durationpb.New(time.Duration(period) * time.Second),
		Depth:  int64(depth),
	}
}
