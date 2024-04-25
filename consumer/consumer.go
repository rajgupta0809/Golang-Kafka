package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "comments"
	consumer, err := connectConsumer([]string{"localhost:29092"})
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer started ")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	msgCount := 0
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Println("err", err)
			case msg := <-partitionConsumer.Messages():
				msgCount++
				fmt.Printf("Received Message Count: %d | Topic(%s) | Message(%s)", msgCount, msg.Topic, string(msg.Value))
			case <-sigchan:
				fmt.Println("Inturption detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed ", msgCount, " messages")
	if err := consumer.Close(); err != nil {
		panic(err)
	}
}

func connectConsumer(brokerUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
