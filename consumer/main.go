package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "coffee_order"
	msgCount := 0
	broker := []string{"localhost:9092"}

	consumer, err := connectConsumer(broker)
	if err != nil {
		panic(err)
	}

	cp, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer")

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-cp.Errors():
				fmt.Println(err)
			case msg := <-cp.Messages():
				msgCount++
				fmt.Printf("Received order Count %d\nTopic %s | Value %s | Timestamp %s\n",
					msgCount, msg.Topic, msg.Value, msg.Timestamp)
			case <-signChan:
				fmt.Println("Interrupt is detected")
				doneChan <- struct{}{}
			}
		}
	}()
	<-doneChan
	fmt.Println("Processed", msgCount, "message")

	if err := consumer.Close(); err != nil {
		panic(err)
	}
}

func connectConsumer(broker []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return sarama.NewConsumer(broker, config)
}
