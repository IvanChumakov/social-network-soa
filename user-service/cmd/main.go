package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"social-network/user-service/internal/app"
	"social-network/user-service/internal/config"
	"social-network/user-service/internal/db"
	"social-network/user-service/internal/logger"
	"social-network/user-service/internal/repository"
	"social-network/user-service/internal/server"
	"social-network/user-service/internal/service"
)

func main() {
	addOpts := fx.Options(
		fx.Provide(
			repository.NewUserRepository,
			service.NewUserService,
			config.NewConfig,
			app.NewApp,
			server.NewServer,
			db.InitDb,
		),
		fx.Invoke(
			godotenv.Load,
			logger.InitLogger,
			server.InvokeServer))
	fx.New(addOpts).Run()
}
