package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	_ "social-network/api-gateway/docs"
	"social-network/api-gateway/internal/app"
	"social-network/api-gateway/internal/client"
	"social-network/api-gateway/internal/config"
	"social-network/api-gateway/internal/logger"
	"social-network/api-gateway/internal/server"
)

// @title Swagger API-GATEWAY
// @version 1.0
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	addOpts := fx.Options(
		fx.Provide(
			config.NewConfig,
			client.NewGrpcConnection,
			client.NewStatsGrpcConnection,
			app.NewApp,
			server.NewServer,
		),
		fx.Invoke(
			godotenv.Load,
			logger.InitLogger,
			server.InvokeServer,
		))
	fx.New(addOpts).Run()
}
