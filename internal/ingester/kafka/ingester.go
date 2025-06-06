package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	kafkarepo "github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"github.com/captainhbb/movieexample-rating/pkg/model"
)

type Ingester struct {
	consumer *kafkarepo.Consumer
	topic    string
}

func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafkarepo.NewConsumer(&kafkarepo.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &Ingester{consumer: consumer, topic: topic}, nil
}

func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
			default:
			}

			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				log.Println("Consumer error: %v", err.Error())
				continue
			}

			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("Unmarshal error: " + err.Error())
				continue
			}
			ch <- event
		}
	}()

	return ch, nil
}
