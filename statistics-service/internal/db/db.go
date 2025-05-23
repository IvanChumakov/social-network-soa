package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"
	"social-network/statistics-service/internal/config"
	"social-network/statistics-service/internal/logger"
)

type Info struct {
	topic    string
	kafkaDB  string
	targetDB string
	mwDB     string
}

func InitDB(lc fx.Lifecycle, cfg config.Config) (*sql.DB, error) {
	conn := sql.OpenDB(clickhouse.Connector(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.DBName,
			Username: cfg.User,
			Password: cfg.Password,
		},
	}))

	if err := conn.Ping(); err != nil {
		logger.Error("error connecting to clickhouse: " + err.Error())
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	meta := []Info{
		{
			topic:    "likes-topic",
			kafkaDB:  "likeskafka",
			targetDB: "likes",
			mwDB:     "likesmw",
		},
		{
			topic:    "views-topic",
			kafkaDB:  "viewskafka",
			targetDB: "views",
			mwDB:     "viewsmw",
		},
		{
			topic:    "comments-topic",
			kafkaDB:  "commentskafka",
			targetDB: "comments",
			mwDB:     "commentsmw",
		},
	}

	for _, info := range meta {
		err := createKafkaTable(conn, info.kafkaDB, info.topic)
		if err != nil {
			return nil, err
		}

		err = createTargetTable(conn, info.targetDB)
		if err != nil {
			return nil, err
		}

		err = createMwTable(conn, info.mwDB, info.targetDB, info.kafkaDB)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func createKafkaTable(conn *sql.DB, name string, topic string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		Time DateTime('UTC'),
		UserId Int32,
		PostId Int32
	)
	ENGINE = Kafka()
	SETTINGS kafka_broker_list = 'kafka:9092',
			 kafka_topic_list = '%s',
			 kafka_group_name = 'clickhouse_%s_consumer',
			 kafka_format = 'JSONEachRow',
			 kafka_num_consumers = 1,
			 kafka_skip_broken_messages = 1,
			 date_time_input_format = 'best_effort'`, name, topic, topic)
	if _, err := conn.Exec(query); err != nil {
		logger.Error(fmt.Sprintf("error creating table %s: %s", name, err.Error()))
		return err
	}

	return nil
}

func createTargetTable(conn *sql.DB, name string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		Time DateTime,
		UserId Int32,
		PostId Int32
	)
	ENGINE = MergeTree()	
	PARTITION BY toYYYYMM(Time)
	ORDER BY (Time)`, name)

	if _, err := conn.Exec(query); err != nil {
		logger.Error(fmt.Sprintf("error creating table %s: %s", name, err.Error()))
		return err
	}

	return nil
}

func createMwTable(conn *sql.DB, name string, to string, from string) error {
	query := fmt.Sprintf(`
			CREATE MATERIALIZED VIEW IF NOT EXISTS %s
			TO %s
			AS 
			SELECT 
				Time,
				UserId,
				PostId
			FROM %s`, name, to, from)
	if _, err := conn.Exec(query); err != nil {
		logger.Error(fmt.Sprintf("error creating table %s: %s", name, err.Error()))
		return err
	}

	return nil
}
