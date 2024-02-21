package fixer

import (
	"context"
	"errors"
	"github.com/Duke1616/etools/kafka/saramax"
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator"
	"github.com/Duke1616/etools/migrator/events"
	"github.com/Duke1616/etools/migrator/fixer"
	"github.com/IBM/sarama"
	"time"
)

type Consumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.Logger
	srcFirst *fixer.OverrideFixer[T]
	dstFirst *fixer.OverrideFixer[T]
	topic    string
}

func NewConsumer[T migrator.Entity](client sarama.Client, l logger.Logger, src *fixer.OverrideFixer[T],
	dst *fixer.OverrideFixer[T], topic string) (*Consumer[T], error) {

	return &Consumer[T]{
		client:   client,
		l:        l,
		srcFirst: src,
		dstFirst: dst,
		topic:    topic,
	}, nil

}
func (c *Consumer[T]) Start() error {
	eg, err := sarama.NewConsumerGroupFromClient("migrator-fix", c.client)
	if err != nil {
		return err
	}
	go func() {
		if err := eg.Consume(context.Background(), []string{c.topic},
			saramax.NewHandler[events.InconsistentEvent](c.Consumer)); err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (c *Consumer[T]) Consumer(msg *sarama.ConsumerMessage, evt events.InconsistentEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch evt.Direction {
	case "SRC":
		return c.srcFirst.Fix(ctx, evt)
	case "DST":
		return c.dstFirst.Fix(ctx, evt)
	}
	return errors.New("未知的校验方向")
}
