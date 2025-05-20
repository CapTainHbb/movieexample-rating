package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/captainhbb/movieexample-rating/pkg/model"
	"os"
	"time"
)

func main() {
	fmt.Println("creating a kafka producer")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const fileName = "ratingsdata.json"
	fmt.Println("reading rating events from file " + fileName)

	ratingEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	const topic = "ratings"
	if err := produceRatingEvents(topic, producer, ratingEvents); err != nil {
		panic(err)
	}

	const timeout = time.Second * 10
	fmt.Printf("waiting %v seconds for rating events to be processed", timeout.String())
	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}

func produceRatingEvents(topic string, producer *kafka.Producer, ratings []model.RatingEvent) error {
	for _, rating := range ratings {
		encodedEvent, err := json.Marshal(rating)
		if err != nil {
			return err
		}

		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(encodedEvent),
		}, nil); err != nil {
			return err
		}
	}
	return nil
}
