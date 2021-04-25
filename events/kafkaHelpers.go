package events

import (
	"example.com/app/config"
	"fmt"
	"github.com/Shopify/sarama"
)

func PushUserToQueue(message []byte) error {

	producer := GetInstance()

	msg := &sarama.ProducerMessage{
		Topic: config.Config("TOPIC"),
		Value: sarama.StringEncoder(message),
	}


	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Failed to send message to the queue")
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "user", partition, offset)
	return nil
}

