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
}

func (s *Suite) SetupSuite() {
	s.log = new(mocks.Logger)
	s.log.On("Info", mock.Anything)
	s.log.On("Error", mock.Anything).Maybe()

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
	collFunc := new(mocks.DataRequestFunc)
	collFunc.On("Execute").Return([]float64{}, nil)
	f := sysmon.Collectors{
		LoadAvg: collFunc.Execute,
		CPU:     collFunc.Execute,
	}

	collector := NewCollector(s.log, s.settings, f)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	err := collector.Start(ctx)
	s.NoError(err)
	time.Sleep(10 * time.Millisecond)
	cancel()
	statusServices := collector.GetStatusServices()
	s.Equal(2, statusServices.La.StatusCode)
	s.Equal(2, statusServices.CPU.StatusCode)
	collFunc.AssertNotCalled(s.T(), "Execute")
}

func (s *Suite) TestServiceTickAndError() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	defer goleak.VerifyNone(s.T())
	collFunc := new(mocks.DataRequestFunc)
	collFunc.On("Execute").Return([]float64{}, errors.New("dummy error"))
	f := sysmon.Collectors{
		LoadAvg: collFunc.Execute,
		CPU:     collFunc.Execute,
	}
	collector := NewCollector(s.log, s.settings, f)

	ctx, cancel := context.WithCancel(context.Background())
	err := collector.Start(ctx)
	s.NoError(err)

	time.Sleep(10 * time.Millisecond)
	result := collector.GetStats(1)
	s.Equal(0, result.La.Counter)
	s.Equal(0, result.CPU.Counter)

	statusServices := collector.GetStatusServices()
	s.Equal(1, statusServices.La.StatusCode)
	s.Equal(1, statusServices.CPU.StatusCode)

	time.Sleep(1 * time.Second)
	result = collector.GetStats(1)
	s.Equal(0, result.La.Counter)
	s.Equal(0, result.CPU.Counter)

	statusServices = collector.GetStatusServices()
	s.Equal(3, statusServices.La.StatusCode)
	s.Equal(3, statusServices.CPU.StatusCode)

	cancel()
	collFunc.AssertExpectations(s.T())
}

func (s *Suite) TestServiceTickDone() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	defer goleak.VerifyNone(s.T())
	collFunc := new(mocks.DataRequestFunc)
	collFunc.On("Execute").Return([]float64{1.0, 2.0, 3.0}, nil)
	f := sysmon.Collectors{
		LoadAvg: collFunc.Execute,
		CPU:     collFunc.Execute,
	}
	collector := NewCollector(s.log, s.settings, f)

	ctx, cancel := context.WithCancel(context.Background())
	err := collector.Start(ctx)
	s.NoError(err)

	time.Sleep(10 * time.Millisecond)
	result := collector.GetStats(1)
	s.Equal(0, result.La.Counter)
	s.Equal(0, result.CPU.Counter)

	statusServices := collector.GetStatusServices()
	s.Equal(1, statusServices.La.StatusCode)
	s.Equal(1, statusServices.CPU.StatusCode)

	time.Sleep(1 * time.Second)
	result = collector.GetStats(1)
	s.Equal(1, result.La.Counter)
	s.Equal(1, result.CPU.Counter)

	statusServices = collector.GetStatusServices()
	s.Equal(1, statusServices.La.StatusCode)
	s.Equal(1, statusServices.CPU.StatusCode)

	cancel()
	collFunc.AssertExpectations(s.T())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
