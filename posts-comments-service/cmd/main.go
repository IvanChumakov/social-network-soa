package main

import (
	"go.uber.org/fx"
	"social-network/posts-comments-service/internal/app"
	"social-network/posts-comments-service/internal/config"
	"social-network/posts-comments-service/internal/db"
	"social-network/posts-comments-service/internal/logger"
	"social-network/posts-comments-service/internal/repository"
	"social-network/posts-comments-service/internal/server"
	eventsservice "social-network/posts-comments-service/internal/service/events-service"
	"social-network/posts-comments-service/internal/service/posts-service"
)

func main() {
	logger.InitLogger()
	addOpts := fx.Options(
		fx.Provide(
			config.NewConfig,
			db.InitDb,
			eventsservice.NewKafkaEvents,
			repository.NewPostRepository,
			func(repo *repository.PostRepository) posts_service.Repository {
				return repo
			},
			posts_service.NewPostService,
			func(service *posts_service.PostService) app.Service {
				return service
			},
			app.NewServer,
		),
		fx.Invoke(
			server.RunServer,
		),
	)
	fx.New(addOpts).Run()
}
