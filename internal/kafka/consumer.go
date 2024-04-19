package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type HandleMessage func(msg []byte) error

type Consumer struct {
	Ready chan bool
	Fn    HandleMessage
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}

			err := c.Fn(message.Value)
			if err == nil {
				session.MarkMessage(message, "")
			} else {
				log.Println(err)
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		case <-session.Context().Done():
			return nil
		}
	}
}
