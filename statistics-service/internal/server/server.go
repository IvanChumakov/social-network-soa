package server

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	pb "social-network/protos"
	"social-network/statistics-service/internal/app"
	"social-network/statistics-service/internal/config"
	"social-network/statistics-service/internal/logger"
)

func RunServer(lc fx.Lifecycle, cfg config.Config, app *app.App) error {
	grpcServer := grpc.NewServer()
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				lis, err := net.Listen("tcp", "statistics-service"+cfg.ServerAddr)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to listen: %v", err))
				}
				logger.Info("starting grpc server on " + cfg.ServerAddr)
				pb.RegisterStatisticsServiceServer(grpcServer, app)

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
