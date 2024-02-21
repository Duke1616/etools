package events

import (
	"github.com/Duke1616/etools/kafka/saramax"
)

type Consumer interface {
	Start() error
}

func InitConsumers(event *saramax.Incr) []Consumer {
	return []Consumer{event}
}
