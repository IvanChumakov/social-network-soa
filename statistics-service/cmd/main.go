package main

import (
	"go.uber.org/fx"
	"social-network/statistics-service/internal/app"
	"social-network/statistics-service/internal/config"
	"social-network/statistics-service/internal/db"
	"social-network/statistics-service/internal/logger"
	"social-network/statistics-service/internal/repository/statistics"
	"social-network/statistics-service/internal/repository/transactor"
	"social-network/statistics-service/internal/server"
	"social-network/statistics-service/internal/service"
)

func main() {
	logger.InitLogger()
	addOpts := fx.Options(
		fx.Provide(
			config.NewConfig,
			db.InitDB,
			statistics.NewStatisticsRepository,
			transactor.NewTxBeginner,
			func(transactor *transactor.TxBeginner) service.Transactor {
				return transactor
			},
			service.NewStatisticsService,
			func(service *service.StatisticsService) app.Service {
				return service
			},
			app.NewApp,
		),
		fx.Invoke(server.RunServer),
	)
	fx.New(addOpts).Run()
}
