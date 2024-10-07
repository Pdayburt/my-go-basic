package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{fn: fn}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			zap.L().Error("Kafka消息反序列化失败",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("offset", msg.Offset),
				zap.Int32("Partition", msg.Partition))
			continue
		}
		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if err == nil {
				break
			}
			zap.L().Error("处理消息失败",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("offset", msg.Offset),
				zap.Int32("Partition", msg.Partition))
		}
		if err != nil {
			zap.L().Error("处理消息失败",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("offset", msg.Offset),
				zap.Int32("Partition", msg.Partition))
		} else {
			session.MarkMessage(msg, "")
		}

	}
	return nil
}
