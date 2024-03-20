package accesslog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type ConsumerService struct {
	kafkaClient   *kgo.Client
	contextCancel context.CancelFunc
}

func NewConsumer(brokers []string, topic string, consumerGroup string) (*ConsumerService, error) {
	log.Printf("Validating configurations...")

	if len(brokers) == 0 {
		return nil, fmt.Errorf("missing brokers")
	}

	if topic == "" {
		return nil, fmt.Errorf("missing topic")
	}

	log.Printf("Initializing kafka consumer...")

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(consumerGroup),
		kgo.ConsumeTopics(topic),
	)

	if err != nil {
		return nil, err
	}

	kcs := &ConsumerService{
		kafkaClient:   cl,
		contextCancel: nil,
	}

	log.Printf("Initialize kafka consumer completed")

	return kcs, nil
}

func (kcs *ConsumerService) Start(process func(accessLog *Log)) error {
	if kcs.contextCancel != nil {
		// Already started, return
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	kcs.contextCancel = cancel

	go func() {
		for {
			fetches := kcs.kafkaClient.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				if len(errs) == 1 && errors.Is(errs[0].Err, context.Canceled) {
					// context stop trigger by normal stop process
					break
				}

				// else
				// All errors are retried internally when fetching, but non-retriable errors are
				// returned from polls so that users can notice and take action.
				kcs.Stop()
				log.Fatalf(fmt.Sprint(errs))
			}

			//log.Printf("A new group of access logs arrived")
			// We can iterate through a callback function.
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				// We can even use a second callback!
				p.EachRecord(func(record *kgo.Record) {
					var accessLog Log
					if err := json.Unmarshal(record.Value, &accessLog); err != nil {
						log.Printf("Failed to parse message into json: %s", string(record.Value))
					}

					process(&accessLog)
				})
			})
		}
	}()

	return nil
}

func (kcs *ConsumerService) Stop() {
	if kcs.contextCancel != nil {
		kcs.contextCancel()
	}

	kcs.kafkaClient.Close()
}
