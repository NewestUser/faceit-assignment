package kfka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaProducer(bootstrapServers string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})

	if err != nil {
		return nil, err
	}

	return &KafkaProducer{p: p}, nil
}

type KafkaProducer struct {
	p *kafka.Producer
}

func (k *KafkaProducer) Produce(topic string, msgValue interface{}) error {
	bts, err := json.Marshal(msgValue)
	if err != nil {
		return fmt.Errorf("failed serializing json kafka message to topic %s err: %s", topic, err)
	}

	err = k.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bts,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed producing kafka message to topic %s err: %s", topic, err)
	}

	return nil
}

func (k *KafkaProducer) Close() {
	k.p.Close()
}
