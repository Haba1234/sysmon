// +build integration

package integration_test

import (
	"context"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"

	stat "github.com/Haba1234/sysmon/internal/grpc/api"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Suite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
	conn   *grpc.ClientConn
	cmd    *exec.Cmd
	port   string
}

func (s *Suite) SetupSuite() {
	var err error

	s.port = os.Getenv("SERVER_PORT")
	if s.port == "" {
		s.port = "8080"
	}

	err = buildServer()
	s.NoError(err)

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.cmd = exec.CommandContext(s.ctx, "./bin/sysmon")
	s.cmd.Dir = "../.."

	err = s.cmd.Start()
	s.NoError(err)
	time.Sleep(500 * time.Millisecond)
	s.conn, err = grpc.DialContext(s.ctx, net.JoinHostPort("", s.port), grpc.WithInsecure())
	s.NoError(err)
}

func (s *Suite) TestClientConnectServer() {
	client := stat.NewStatisticsClient(s.conn)
	stream := s.runClient(client, 1, 1)
	stats, err := stream.Recv()
	s.NoError(err)
	s.Equal("works", stats.La.GetStatus())
	s.Equal("works", stats.Cp.GetStatus())
}

func (s *Suite) TestLoadCPU() {
	client := stat.NewStatisticsClient(s.conn)
	stream := s.runClient(client, 1, 1)
	stats, err := stream.Recv()
	s.NoError(err)
	loadCPUStart, _ := strconv.ParseFloat(stats.Cp.GetUser(), 64)

	go loadCPU(2 * time.Second)
	stats, err = stream.Recv()
	s.NoError(err)
	loadCPUStop, _ := strconv.ParseFloat(stats.Cp.GetUser(), 64)
	s.Greater(loadCPUStop, loadCPUStart)
}

func (s *Suite) TestManyClients() {
	wg := sync.WaitGroup{}
	const count = 10
	wg.Add(count - 1)

	for i := 1; i < count; i++ {
		go func(m int) {
			defer wg.Done()
			client := stat.NewStatisticsClient(s.conn)
			stream := s.runClient(client, m, 1)
			_, err := stream.Recv()
			s.NoError(err)
		}(i)
	}
	wg.Wait()
}

func (s *Suite) TearDownSuite() {
	s.conn.Close()
	s.cancel()
	_ = s.cmd.Wait()
	err := os.Remove("../../bin/sysmon")
	s.NoError(err)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func buildServer() error {
	cmd := exec.Command("make", "build")
	cmd.Dir = "../.."
	return cmd.Run()
}

func (s *Suite) runClient(statClient stat.StatisticsClient, m, n int) stat.Statistics_ListStatisticsClient {
	req := &stat.SubscriptionRequest{
		Period: durationpb.New(time.Duration(n) * time.Second),
		Depth:  int64(m),
	}

	stream, err := statClient.ListStatistics(s.ctx, req)
	s.NoError(err)

	return stream
}

func loadCPU(d time.Duration) {
	begin := time.Now()
	for {
		if time.Since(begin) > d {
			break
		}
	}
}
