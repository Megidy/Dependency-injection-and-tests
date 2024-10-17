package producer

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewProducer(producer sarama.SyncProducer) *KafkaProducer {
	return &KafkaProducer{
		producer: producer,
	}

}

func (p *KafkaProducer) PushOrderToQueue(topic, producerPort string, message []byte) error {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := p.producer.SendMessage(&msg)
	if err != nil {
		return err
	}
	log.Printf("Order Message is stored in topic(%s)/partition(%d)/offset(%d)",
		topic, partition, offset)
	return nil
}
