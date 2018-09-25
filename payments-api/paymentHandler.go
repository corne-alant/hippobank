package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//PaymentResponse DTO
type PaymentResponse struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func writeToKafka() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
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

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries
	p.Flush(15 * 1000)
}

// simpleHandler returns current date time
func (s *server) paymentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		payment := PaymentResponse{
			From:   "Account 1",
			To:     "Account 2",
			Amount: 200,
		}
		resJSON, err := json.Marshal(payment)
		if err != nil {
			panic(err)
		}

		writeToKafka()

		//Set Content-Type header so that clients will know how to read response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//Write json response back to response
		w.Write(resJSON)
	}
}
