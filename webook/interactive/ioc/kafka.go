package ioc

import (
	"example.com/mod/webook/interactive/events"
	"example.com/mod/webook/pkg/saramax"
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

// NewConsumer 这个模块没哟Kafka的producer
/*func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	syncProducer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return syncProducer
}
*/
// NewConsumer 多个消费者在这里注册
func NewConsumer(cl *events.InteractiveReadEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{cl}
}
