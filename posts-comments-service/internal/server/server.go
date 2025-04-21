package server

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	"social-network/posts-comments-service/internal/app"
	"social-network/posts-comments-service/internal/config"
	"social-network/posts-comments-service/internal/logger"
	pb "social-network/protos"
)

func RunServer(lc fx.Lifecycle, cfg *config.Config, server *app.Server) error {
	grpcServer := grpc.NewServer()
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				lis, err := net.Listen("tcp", "posts-service"+cfg.ServAddr)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to listen: %v", err))
				}
				logger.Info("starting grpc server on " + cfg.ServAddr)
				pb.RegisterPostsServiceServer(grpcServer, server)

				if err = grpcServer.Serve(lis); err != nil {
					logger.Error("error starting server: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
	return nil
}
