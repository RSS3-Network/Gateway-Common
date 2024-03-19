package accesslog_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/rss3-network/gateway-common/accesslog"
)

func TestAccessLog(t *testing.T) {
	t.Parallel()

	// Prepare configs
	brokers := []string{"localhost:19092"}
	topic := "gateway.common.test"
	consumerGroup := "gateway-common-test"

	// Create producer
	producer, err := accesslog.NewProducer(brokers, topic)

	if err != nil {
		t.Fatal(fmt.Errorf("create producer: %w", err))
	}

	defer producer.Stop()

	// Create consumer
	consumer, err := accesslog.NewConsumer(brokers, topic, consumerGroup)

	if err != nil {
		t.Fatal(fmt.Errorf("create consumer: %w", err))
	}

	defer consumer.Stop()

	// Prepare test case storage space
	demoLogs := []accesslog.Log{
		{
			Key:       nil, // No key
			Path:      "/foo",
			Status:    http.StatusOK,
			Timestamp: time.Unix(1710849419, 0),
		},
		{
			Key:       toPtr("84b01bc1-4dad-4694-99ce-514c37b88f9a"),
			Path:      "/bar",
			Status:    http.StatusTooManyRequests,
			Timestamp: time.Unix(1710849621, 0),
		},
		{
			Key:       nil, // No key
			Path:      "/baz",
			Status:    http.StatusInternalServerError,
			Timestamp: time.Unix(1710849652, 0),
		},
		{
			Key:       toPtr("7eeb2c6d-d94f-475b-907c-50cbe01a0cb6"),
			Path:      "/bar?alice=bob",
			Status:    http.StatusTooManyRequests,
			Timestamp: time.Unix(1710849711, 0),
		},
	}
	receiveLogChan := make(chan accesslog.Log, len(demoLogs)+1)

	var wg sync.WaitGroup

	// Start consuming
	if err = consumer.Start(func(accessLog *accesslog.Log) {
		t.Log("access log consume")
		receiveLogChan <- *accessLog
		wg.Done()
	}); err != nil {
		t.Error(err)
	}

	// Start producing
	for _, l := range demoLogs {
		t.Log("access log produce")

		l := l

		if err = producer.ProduceLog(&l); err != nil {
			t.Error(err)
		}

		wg.Add(1)
	}

	// Wait for all process finish
	t.Log("waiting all finish...")

	wg.Wait()

	// Close channel
	close(receiveLogChan)

	// Compare results
	counter := 0

	for receivedLog := range receiveLogChan {
		t.Log(receivedLog, demoLogs[counter])

		if !((receivedLog.Key == nil && demoLogs[counter].Key == nil) ||
			*receivedLog.Key == *demoLogs[counter].Key) {
			t.Error(fmt.Errorf("item %d key mismatch", counter))
		}

		if receivedLog.Path != demoLogs[counter].Path {
			t.Error(fmt.Errorf("item %d path mismatch", counter))
		}

		if receivedLog.Status != demoLogs[counter].Status {
			t.Error(fmt.Errorf("item %d status mismatch", counter))
		}

		if receivedLog.Timestamp.UnixNano() != demoLogs[counter].Timestamp.UnixNano() {
			t.Error(fmt.Errorf("item %d ts mismatch", counter))
		}

		counter++

		if counter > len(demoLogs) {
			break
		}
	}

	if counter != len(demoLogs) {
		t.Error("invalid logs length")
	}

	t.Log("test finish")
}

func toPtr[T any](v T) *T {
	return &v
}
