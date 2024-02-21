package startup

import (
	"github.com/Duke1616/etools/logger"
	"github.com/Duke1616/etools/migrator/events"
	"github.com/Duke1616/etools/migrator/events/fixer"
	"github.com/Duke1616/etools/migrator/example"
	"github.com/IBM/sarama"
)

func InitFixerProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer("inconsistent_users", p)
}

func InitFixerConsumer(client sarama.Client, l logger.Logger, src SrcDB,
	dst DstDB) *fixer.Consumer[example.User] {
	res, err := fixer.NewConsumer[example.User](client, l, src, dst, "inconsistent_users")
	if err != nil {
		panic(err)
	}
	return res
}
