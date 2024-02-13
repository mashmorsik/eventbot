package cronkafka

import (
	"encoding/json"
	"eventbot/data"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

const (
	broker = "localhost:19092"
	topic  = "table-kafka"
)

func Consumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "events-group",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}

	msg, err := c.ReadMessage(time.Second * 3)
	if err != nil {
		fmt.Println(msg)
	}
	//select {
	//case <-time.After(time.Second * 3):
	//	fmt.Println("No news in an hour. Service aborting.")
	//default:
	//	msg, err := c.ReadMessage(time.Second)
	//	if err == nil {
	//		fmt.Printf("Consumed message from %s: %s\n", msg.TopicPartition, string(msg.Value))
	//	} else if !err.(kafka.Error).IsTimeout() {
	//		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	//	}
	//}

	c.Close()
}

func Producer(event *data.Event) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	eventByte, err := json.Marshal(event)

	// Produce messages to topic (asynchronously)
	topic := topic
	err = p.Produce(&kafka.Message{
		Key:            []byte("NewEvent"),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          eventByte,
	}, nil)
	if err != nil {
		panic(err)
	}
	//p.Flush()
}
