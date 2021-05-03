package main

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	stat "github.com/Haba1234/sysmon/internal/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client start error: %v", err)
	}
	defer conn.Close()

	client := stat.NewStatisticsClient(conn)
	createClient(client)
}

func createClient(statClient stat.StatisticsClient) {
	req := &stat.SubscriptionRequest{
		Period: durationpb.New(5 * time.Second),
		Depth:  5,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	stream, err := statClient.ListStatistics(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Internal {
			log.Println("Внутренняя ошибка сервера")
		} else {
			log.Println("cannot create client gRPC: ", err)
		}
		return
	}

	for {
		stats, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("Получен конец передачи данных, error: ", err)
			return
		}
		if err != nil {
			log.Println("Поток данных прерван, error: ", err)
			return
		}
		log.Printf("LoadAverage: %v", stats.GetLa())
	}
}
