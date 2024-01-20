package example

import (
	"encoding/json"
	"github.com/Duke1616/etools/kafka/sarama/event"
	"github.com/Duke1616/etools/kafka/sarama/saramax"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const EvtTopic = "incr_topic"

func TestHandler(t *testing.T) {
	// 客户端
	client := saramax.SaramaClient()

	// 消费者
	incr := saramax.NewIncr(client, EvtTopic)
	consumer := event.InitConsumers(incr)

	for _, c := range consumer {
		er := c.Start()
		if er != nil {
			panic(er)
		}
	}

	time.Sleep(100 * time.Second)
}

func TestAsyncProducers(t *testing.T) {
	client := saramax.SaramaClient()
	producer, err := sarama.NewSyncProducerFromClient(client)
	assert.NoError(t, err)

	// 序列号Event
	evt := saramax.EventIncr{Name: "你好我是tools"}
	val, err := json.Marshal(evt)
	assert.NoError(t, err)

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: EvtTopic,
		Value: sarama.StringEncoder(val),
	})

	assert.NoError(t, err)
}
