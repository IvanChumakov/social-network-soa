package main

import (
	"go.uber.org/fx"
	"social-network/posts-comments-service/internal/app"
	"social-network/posts-comments-service/internal/config"
	"social-network/posts-comments-service/internal/db"
	"social-network/posts-comments-service/internal/logger"
	"social-network/posts-comments-service/internal/repository"
	"social-network/posts-comments-service/internal/server"
	"social-network/posts-comments-service/internal/service"
)

func main() {
	logger.InitLogger()
	addOpts := fx.Options(
		fx.Provide(
			config.NewConfig,
			db.InitDb,
			repository.NewPostRepository,
			func(repo *repository.PostRepository) service.Repository {
				return repo
			},
			service.NewPostService,
			func(service *service.PostService) app.Service {
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
