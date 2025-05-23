package statistics

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/statistics-service/internal/logger"
	"social-network/statistics-service/internal/repository/transactor"
	"time"
)

type StatisticsRepository struct {
	db *sql.DB
}

func NewStatisticsRepository(db *sql.DB) *StatisticsRepository {
	return &StatisticsRepository{
		db: db,
	}
}

func (sr *StatisticsRepository) GetLikesCount(ctx context.Context, postId int32) (int32, error) {
	querier := transactor.GetQueries(ctx, sr.db)
	query := `SELECT COUNT(*) FROM likes WHERE PostId = ?`
	var count int32

	if err := querier.QueryRow(query, postId).Scan(&count); err != nil {
		logger.Error("error getting likes count: " + err.Error())
		return 0, err
	}

	return count, nil
}

func (sr *StatisticsRepository) GetViewsCount(ctx context.Context, postId int32) (int32, error) {
	querier := transactor.GetQueries(ctx, sr.db)
	query := `SELECT COUNT(*) FROM views WHERE PostId = ?`
	var count int32

	if err := querier.QueryRow(query, postId).Scan(&count); err != nil {
		logger.Error("error getting views count: " + err.Error())
		return 0, err
	}

	return count, nil
}

func (sr *StatisticsRepository) GetCommentsCount(ctx context.Context, postId int32) (int32, error) {
	querier := transactor.GetQueries(ctx, sr.db)
	query := `SELECT COUNT(*) FROM comments WHERE PostId = ?`
	var count int32

	if err := querier.QueryRow(query, postId).Scan(&count); err != nil {
		logger.Error("error getting comments count: " + err.Error())
		return 0, err
	}

	return count, nil
}

type Dynamic struct {
	Count int32
	Time  time.Time
}

func (sr *StatisticsRepository) GetDynamic(ctx context.Context, postId int32, table string) ([]Dynamic, error) {
	querier := transactor.GetQueries(ctx, sr.db)
	query := fmt.Sprintf(`SELECT 
				toDate(Time) AS date,
				count() AS view_count
			FROM 
				%s
			WHERE 
				PostId = ?
			GROUP BY 
				date
			ORDER BY 
				date ASC`, table)
	rows, err := querier.Query(query, postId)
	if err != nil {
		logger.Error("error getting views dynamic count: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	dynamics := make([]Dynamic, 0)

	for rows.Next() {
		var dynamic Dynamic
		if err = rows.Scan(&dynamic.Time, &dynamic.Count); err != nil {
			logger.Error("error scanning rows for views dynamic: " + err.Error())
			return nil, err
		}

		dynamics = append(dynamics, dynamic)
	}

	return dynamics, nil
}

func (sr *StatisticsRepository) GetTop(ctx context.Context, table string, idType string) ([]int32, error) {
	querier := transactor.GetQueries(ctx, sr.db)
	query := fmt.Sprintf(`SELECT 
				%s,
				count() AS total_views
			FROM 
				%s
			GROUP BY 
				%s
			ORDER BY 
				total_views DESC
			LIMIT 10`, idType, table, idType)
	rows, err := querier.Query(query)
	if err != nil {
		logger.Error(fmt.Sprintf("error getting top %s count: ", table) + err.Error())
		return nil, err
	}
	defer rows.Close()

	ids := make([]int32, 0)
	for rows.Next() {
		var ignore sql.RawBytes
		var id int32
		if err = rows.Scan(&id, &ignore); err != nil {
			logger.Error("error scanning rows for top views: " + err.Error())
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
