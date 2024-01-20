package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {
	return &Handler[T]{
		fn: fn,
	}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for msg := range messages {
		// 在这里调用业务处理逻辑
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			// 你也可以在这里引入重试的逻辑
			slog.Default().Error("反序列消息体失败",
				slog.String("topic", msg.Topic),
				slog.Int64("partition", int64(msg.Partition)),
				slog.Int64("offset", msg.Offset),
				slog.Any("error", err))
		}
		err = h.fn(msg, t)
		if err != nil {
			slog.Default().Error("处理消息失败",
				slog.String("topic", msg.Topic),
				slog.Int64("partition", int64(msg.Partition)),
				slog.Int64("offset", msg.Offset),
				slog.Any("error", err))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
