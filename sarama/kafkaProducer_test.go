package sarama

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncKafkaProducer(t *testing.T) {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("~~~~~~~~~hello Go Kafka~~~~~~"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "这个metadata",
	})
	assert.NoError(t, err)
}

func TestSyncKafkaProducerSendReadEvent(t *testing.T) {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(ReadEvent{
		Uid: 987,
		Aid: 123,
	})
	assert.NoError(t, err)

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "article_read",
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "这个metadata",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)
	msgCh := producer.Input()
	msgCh <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("~~~~~~ Async~~~hello Go Kafka~~~~~~"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "这个metadata",
	}
	errCh := producer.Errors()
	seccCh := producer.Successes()

	select {
	case err := <-errCh:
		t.Log("发送失败", err.Err, err.Msg)
	case suc := <-seccCh:
		t.Log("发送成功", suc.Topic, suc.Key, suc.Value)
	}
}

type ReadEvent struct {
	Uid int64
	Aid int64
}
