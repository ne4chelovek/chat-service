package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"strings"
)

type Handler interface {
	HandleMessage(message []byte, offset kafka.Offset) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(handler Handler, address []string, topic, consumerGroup string) (*Consumer, error) {
	cons, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       6000, //ms
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  6000,
		"auto.offset.reset":        "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = cons.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: cons, handler: handler}, nil
}

func (c *Consumer) Start() {
	for range 1 {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error reading message from consumer: %v", err)
		}
		if kafkaMsg == nil {
			continue
		}
		if err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset); err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}
		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			log.Printf("Error storing message: %v", err)
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}

	log.Printf("Commited offset")
	return c.consumer.Close()
}
