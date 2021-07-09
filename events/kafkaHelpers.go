package events

import (
	"example.com/app/config"
	"example.com/app/domain"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func PushUserToQueue(message []byte) error {

	producer := GetInstance()

	msg := &sarama.ProducerMessage{
		Topic: config.Config("TOPIC"),
		Value: sarama.StringEncoder(message),
	}


	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		err = producer.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("Failed to send message to the queue")
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", "user", partition, offset)
	return nil
}

func SendKafkaMessage(user *domain.User, eventType int) error {
	um := new(domain.UserMessage)
	um.User = user
	fmt.Println(um.User)

	// user created/updated event
	um.MessageType = eventType

	// turn user struct into a byte array
	b, err := msgpack.Marshal(&um)

	if err != nil {
		return err
	}

	err = PushUserToQueue(b)

	if err != nil {
		return err
	}

	return nil
}

func HandleKafkaMessage(err error, user *domain.User, messageType int) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	err = SendKafkaMessage(user, messageType)

	if err != nil {
		fmt.Println("Failed to publish new user")
	}

	return nil
}