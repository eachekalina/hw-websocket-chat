package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"homework/internal/kafka"
	"homework/internal/model"
	"homework/internal/pb"
)

type MessageRepository struct {
	client   pb.MessageStorageClient
	producer *kafka.Producer
}

func NewMessageRepository(client pb.MessageStorageClient, producer *kafka.Producer) *MessageRepository {
	return &MessageRepository{client: client, producer: producer}
}

func (r *MessageRepository) AddMessage(ctx context.Context, message model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json: %w", err)
	}
	err = r.producer.SendMessage(data)
	if err != nil {
		return fmt.Errorf("kafka: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetLastMessages(ctx context.Context, n int) ([]model.Message, error) {
	req := &pb.GetLastMessagesRequest{Number: int32(n)}
	msgs, err := r.client.GetLastMessages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc: %w", err)
	}
	res := make([]model.Message, len(msgs.Messages))
	for i, msg := range msgs.Messages {
		res[i] = model.Message{
			SentTime: msg.SentTime.AsTime(),
			Nickname: msg.Nickname,
			Message:  msg.Message,
		}
	}
	return res, nil
}
