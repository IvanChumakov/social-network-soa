package client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"social-network/api-gateway/internal/config"
	"social-network/api-gateway/internal/logger"
	pb "social-network/protos"
	"syscall"
)

func NewGrpcConnection(cfg *config.Config) (pb.PostsServiceClient, error) {
	conn, err := grpc.NewClient("posts-service"+cfg.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("error connecting to grpc server: %v", err))
		return nil, err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		conn.Close()
	}()

	return pb.NewPostsServiceClient(conn), nil
}

func NewStatsGrpcConnection(cfg *config.Config) (pb.StatisticsServiceClient, error) {
	conn, err := grpc.NewClient("statistics-service"+cfg.StatsGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("error connecting to stats grpc server: %v", err))
		return nil, err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		conn.Close()
	}()

	return pb.NewStatisticsServiceClient(conn), nil
}
