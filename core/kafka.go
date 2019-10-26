package core

import (
	"fmt"

	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
)

// KafkaTransport implements transport interface for kafka protocol
type KafkaTransport struct {
	producer sarama.SyncProducer
}

const (
	defaultMaxRetry = 5
)

// Init creates a kafka producer
func (k *KafkaTransport) Init(url string, token string) error {
	// generates random uuid based on timestamp and mac address
	id := uuid.NewV1()

	config := sarama.NewConfig()
	config.ClientID = fmt.Sprintf("I1820-lg-%v", id)
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = defaultMaxRetry

	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll

	sp, err := sarama.NewSyncProducer([]string{url}, config)
	if err != nil {
		return err
	}

	k.producer = sp

	return nil
}

// Transmit sends data on given kafka topic
func (k *KafkaTransport) Transmit(topic string, data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(data),
	}

	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
