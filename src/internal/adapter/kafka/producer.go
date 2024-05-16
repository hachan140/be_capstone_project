package kafka

import (
	"be-capstone-project/src/internal/core/events"
	"be-capstone-project/src/internal/core/logger"
	"be-capstone-project/src/internal/core/utils"
	"context"
	"github.com/Shopify/sarama"
)

type IEventService interface {
	SendMessage(c context.Context, topic string, key string, event events.Event) error
}

type KafkaService struct {
	SyncProducer sarama.SyncProducer
}

func NewSyncProducerKafkaService(addr []string, config *sarama.Config) (IEventService, error) {
	syncProducer, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		return nil, err
	}

	kafkaService := KafkaService{
		SyncProducer: syncProducer,
	}

	return &kafkaService, nil
}

func (k *KafkaService) SendMessage(c context.Context, topic string, key string, event events.Event) error {
	eventByte, err := utils.EncodeEvent(event)
	if err != nil {
		return err
	}

	message := string(eventByte)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := k.SyncProducer.SendMessage(msg)

	if err != nil {
		return err
	}

	logger.InfoCtx(c, "Success produce message to topic %v, user_id %v, partition %v offset %v message %v",
		topic, event.UserID, partition, offset, message)
	return nil
}
