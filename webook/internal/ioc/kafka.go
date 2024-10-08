package ioc

import (
	events2 "example.com/mod/webook/interactive/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaConfig)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	syncProducer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return syncProducer
}

// 多个消费者在这里注册
func NewConsumer(cl *events2.InteractiveReadEventConsumer) []events2.Consumer {
	return []events2.Consumer{cl}

}
