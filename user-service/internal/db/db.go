package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"social-network/user-service/internal/config"
	"social-network/user-service/internal/logger"
	"social-network/user-service/internal/repository"
)

func InitDb(cfg *config.Config) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@postgres:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresPort, cfg.PostgresDb)

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("init db err: %s", err.Error()))
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*repository.User)(nil)).
		Exec(context.Background())

	if err != nil {
		logger.Error(fmt.Sprintf("create table err: %s", err.Error()))
	}
	logger.Info(fmt.Sprintf("init db success: %s", dsn))
	return db
}
