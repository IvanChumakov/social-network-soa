package events_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"hash/adler32"
	"social-network/posts-comments-service/internal/config"
	"social-network/posts-comments-service/internal/logger"
)

type KafkaEvents struct {
	writer *kafka.Writer
}

func NewKafkaEvents(cfg *config.Config, lc fx.Lifecycle) (*KafkaEvents, error) {
	client := kafka.Client{
		Addr: kafka.TCP(cfg.KafkaUrl + cfg.KafkaPort),
	}

	writer := kafka.Writer{
		Addr:      kafka.TCP(cfg.KafkaUrl + cfg.KafkaPort),
		Balancer:  &kafka.Hash{Hasher: adler32.New()},
		Transport: kafka.DefaultTransport,
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return writer.Close()
		},
	})

	createTopicsReq := kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{
			{
				Topic:             "views-topic",
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
			{
				Topic:             "comments-topic",
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
			{
				Topic:             "likes-topic",
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		},
	}

	createTopic, err := client.CreateTopics(context.Background(), &createTopicsReq)
	if err != nil {
		logger.Error(fmt.Sprintf("error creating topics: %s", err.Error()))
		return nil, err
	}

	for topic, err := range createTopic.Errors {
		if err != nil {
			logger.Error(fmt.Sprintf("error creating topic %s: %s", topic, err.Error()))
		}
	}

	return &KafkaEvents{
		writer: &writer,
	}, nil
}

func (ke *KafkaEvents) SendEvent(eventModel any, topic string) error {
	err := baseProducer{
		produceWithLog{
			ke.writer,
		},
	}.SendMessage(context.Background(), eventModel, topic)
	if err != nil {
		logger.Error(fmt.Sprintf("error sending event: %s", err.Error()))
		return err
	}

	return nil
}

type producer interface {
	WriteMessages(context.Context, ...kafka.Message) error
}

type baseProducer struct {
	producer
}

type produceWithLog struct {
	producer
}

func (p produceWithLog) WriteMessages(ctx context.Context, messages ...kafka.Message) error {
	logger.Info("preparing to send event")

	err := p.producer.WriteMessages(ctx, messages...)
	if err != nil {
		logger.Error(fmt.Sprintf("error sending event: %v", err))
		return err
	}

	logger.Info("notification sent")

	return nil
}

func (bp baseProducer) SendMessage(ctx context.Context, msg any, topic string) error {
	byteData, err := json.Marshal(msg)
	if err != nil {
		logger.Error("error marshaling data")
		return err
	}

	kafkaMsg := kafka.Message{
		Topic: topic,
		Value: byteData,
	}

	err = bp.producer.WriteMessages(ctx, kafkaMsg)

	return err
}
