package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"social-network/posts-comments-service/internal/config"
	"social-network/posts-comments-service/internal/logger"
	"social-network/posts-comments-service/internal/repository"
)

func InitDb(cfg *config.Config) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@posts-postgres:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresPort, cfg.PostgresDb)

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("init db err: %s", err.Error()))
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*repository.Comment)(nil)).
		Exec(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("create comment table err: %s", err.Error()))
		return nil
	}

	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*repository.Post)(nil)).
		Exec(context.Background())

	if err != nil {
		logger.Error(fmt.Sprintf("create table err: %s", err.Error()))
	}
	logger.Info(fmt.Sprintf("init db success: %s", dsn))
	return db
}
