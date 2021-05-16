package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Haba1234/sysmon/internal/mocks"
	"github.com/Haba1234/sysmon/internal/sysmon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type Suite struct {
	suite.Suite
	log      *mocks.Logger
	settings Collection
	m        int
}

func (s *Suite) SetupSuite() {
	s.log = new(mocks.Logger)
	s.log.On("Info", mock.Anything)
	s.log.On("Error", mock.Anything).Maybe()
	s.m = 1
	s.settings = Collection{
		LoadAverageEnabled: true,
		CPUEnabled:         true,
		BufSize:            10,
	}
}

func (s *Suite) TearDownTest() {
	s.log.AssertExpectations(s.T())
}

func (s *Suite) TestServiceCancel() {
	defer goleak.VerifyNone(s.T())
	collector, collFunc := s.createCollector([]float64{}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	err := collector.Start(ctx)
	s.NoError(err)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(10 * time.Millisecond)
	s.getStatusEqual(collector, sysmon.ServiceStop)

	collFunc.AssertNotCalled(s.T(), "Execute")
}

func (s *Suite) TestServiceTickAndError() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	defer goleak.VerifyNone(s.T())
	collector, collFunc := s.createCollector([]float64{}, errors.New("dummy error"))

	ctx, cancel := context.WithCancel(context.Background())
	err := collector.Start(ctx)
	s.NoError(err)

	time.Sleep(10 * time.Millisecond)
	s.getStatEqual(collector, s.m, 0)
	s.getStatusEqual(collector, sysmon.ServiceRun)

	time.Sleep(1200 * time.Millisecond)
	s.getStatEqual(collector, s.m, 0)
	s.getStatusEqual(collector, sysmon.ServiceError)

	cancel()
	collFunc.AssertExpectations(s.T())
}

func (s *Suite) TestServiceTickDone() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	defer goleak.VerifyNone(s.T())
	collector, collFunc := s.createCollector([]float64{1.0, 2.0, 3.0}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	err := collector.Start(ctx)
	s.NoError(err)

	time.Sleep(10 * time.Millisecond)
	s.getStatEqual(collector, s.m, 0)
	s.getStatusEqual(collector, sysmon.ServiceRun)

	time.Sleep(1200 * time.Millisecond)
	s.getStatEqual(collector, s.m, 1)
	s.getStatusEqual(collector, sysmon.ServiceRun)

	cancel()
	collFunc.AssertExpectations(s.T())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) createCollector(data []float64, err error) (sysmon.Collector, *mocks.DataRequestFunc) {
	collFunc := new(mocks.DataRequestFunc)
	collFunc.On("Execute").Return(data, err)
	f := sysmon.Collectors{
		LoadAvg: collFunc.Execute,
		CPU:     collFunc.Execute,
	}
	return NewCollector(s.log, s.settings, f), collFunc
}

func (s *Suite) getStatEqual(c sysmon.Collector, m, expected int) {
	result := c.GetStats(m)
	s.Equal(expected, result.La.Counter)
	s.Equal(expected, result.CPU.Counter)
}

func (s *Suite) getStatusEqual(c sysmon.Collector, expected int) {
	status := c.GetStatusServices()
	s.Equal(expected, status.La.StatusCode)
	s.Equal(expected, status.CPU.StatusCode)
}
