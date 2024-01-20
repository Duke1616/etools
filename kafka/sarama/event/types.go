package event

import "github.com/Duke1616/etools/kafka/sarama/saramax"

type Consumer interface {
	Start() error
}

func InitConsumers(event *saramax.Incr) []Consumer {
	return []Consumer{event}
}
