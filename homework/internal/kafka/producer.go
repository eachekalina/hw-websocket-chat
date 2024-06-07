package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
)

type Producer struct {
	prod  sarama.SyncProducer
	topic string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	prod, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("kafka producer: %w", err)
	}

	return &Producer{prod: prod, topic: topic}, nil
}

func (p *Producer) SendMessage(data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(data),
	}
	_, _, err := p.prod.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka: %w", err)
	}
	return nil
}
