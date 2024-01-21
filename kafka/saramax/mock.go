package saramax

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"time"
)

type Incr struct {
	client sarama.Client
	topic  []string
}

func NewIncr(client sarama.Client, topic string) *Incr {
	return &Incr{
		client: client,
		topic:  []string{topic},
	}
}

func SaramaClient() sarama.Client {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9094"}, conf)
	if err != nil {
		panic(err)
	}
	return client
}

func (i *Incr) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("incr", i.client)
	if err != nil {
		return err
	}

	go func() {
		er := cg.Consume(context.Background(), i.topic, NewHandler[EventIncr](i.Consume))
		if er != nil {
			slog.Default().Error("退出消费", er)
		}
	}()

	return err
}

type EventIncr struct {
	Name string
}

func (i *Incr) Consume(msg *sarama.ConsumerMessage,
	event EventIncr) error {
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println(event.Name)
	return nil
}
