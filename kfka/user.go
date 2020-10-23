package kfka

import (
	"github.com/newestuser/faceit/user"
)

const userTopic = "user"

func NewKafkaUserEventEmitter(k *KafkaProducer) user.EventEmitter {
	return &kafkaUserEventEmitter{producer: k}
}

type kafkaUserEventEmitter struct {
	producer *KafkaProducer
}

func (k kafkaUserEventEmitter) EmitUpdate(e user.Event, u *user.User) error {
	return k.producer.Produce(userTopic, &kafkaUserEvent{Typ: string(e), Data: u})
}

type kafkaUserEvent struct {
	Typ  string     `json:"type"`
	Data *user.User `json:"data"`
}
