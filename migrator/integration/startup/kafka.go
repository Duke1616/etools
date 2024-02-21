package startup

import (
	"github.com/Duke1616/etools/kafka/saramax/events"
	"github.com/Duke1616/etools/migrator/events/fixer"
	"github.com/Duke1616/etools/migrator/example"
	"github.com/IBM/sarama"
)

func InitSaramaClient() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9094"}, cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}

func InitConsumers(fixConsumer *fixer.Consumer[example.User]) []events.Consumer {
	return []events.Consumer{fixConsumer}
}
