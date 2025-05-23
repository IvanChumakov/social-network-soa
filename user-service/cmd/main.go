package main

import (
	"social-network/user-service/internal/app"
	"social-network/user-service/internal/config"
	"social-network/user-service/internal/db"
	"social-network/user-service/internal/logger"
	"social-network/user-service/internal/repository"
	"social-network/user-service/internal/server"
	eventservice "social-network/user-service/internal/service/event-service"
	user_service "social-network/user-service/internal/service/user-service"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	addOpts := fx.Options(
		fx.Provide(
			eventservice.NewKafkaEvents,
			repository.NewUserRepository,
			user_service.NewUserService,
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
