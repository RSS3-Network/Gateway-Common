package accesslog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type ProducerClient struct {
	kafkaClient *kgo.Client
	topic       string
}

func NewProducer(brokers []string, topic string) (*ProducerClient, error) {
	log.Printf("Validating configurations...")

	if len(brokers) == 0 {
		return nil, fmt.Errorf("missing brokers")
	}

	if topic == "" {
		return nil, fmt.Errorf("missing topic")
	}

	log.Printf("Initializing kafka client...")

	cl, err := kgo.NewClient(
		kgo.AllowAutoTopicCreation(),
		kgo.SeedBrokers(brokers...),
	)

	if err != nil {
		return nil, err
	}

	kcs := &ProducerClient{
		kafkaClient: cl,
		topic:       topic,
	}

	log.Printf("Initialize kafka producer completed")

	return kcs, nil
}

func (kpc *ProducerClient) ProduceLog(l *Log) error {
	logBytes, err := json.Marshal(&l)
	if err != nil {
		return fmt.Errorf("marshal bytes %v with error: %w", l, err)
	}

	record := &kgo.Record{
		Topic: kpc.topic,
		Value: logBytes,
	}

	kpc.kafkaClient.Produce(context.Background(), record, func(_ *kgo.Record, err error) {
		if err != nil {
			log.Printf("Failed to produce record %v: %v", record, err)
		}
	})

	return nil
}

func (kpc *ProducerClient) Stop() {
	kpc.kafkaClient.Close()
}
