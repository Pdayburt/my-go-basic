package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {

	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()

	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancelFunc()

	ctx, cancelFunc := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*3, func() {
		cancelFunc()
	})

	err = consumer.Consume(ctx,
		[]string{"test_topic"}, testConsumerGroupHandler{})
	t.Log(err, time.Since(start).String())

}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {

	//topic 偏移量
	//每次重启setup都会执行
	partions := session.Claims()["test_topic"]
	for _, partion := range partions {
		session.ResetOffset("test_topic", partion,
			sarama.OffsetOldest, "")
	}

	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {

	log.Print("testing cleanup")
	return nil
}

func (t testConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		var eg errgroup.Group
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		msgTOConsumer := make([]*sarama.ConsumerMessage, 0, 10)
		for i := 0; i < batchSize; i++ {
			select {
			case <-ctx.Done():
				cancelFunc()
			case msg, ok := <-msgs:
				if !ok {
					cancelFunc()
				}

				msgTOConsumer = append(msgTOConsumer, msg)
				eg.Go(func() error {
					log.Println(msg)
					return nil
				})

			}
		}
		err := eg.Wait()
		if err != nil {
			log.Print(err)
		}

		for _, message := range msgTOConsumer {
			session.MarkMessage(message, "")
		}

	}
}

// 同步消费
func (t testConsumerGroupHandler) ConsumeClaimV1(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for val := range messages {
		log.Println(string(val.Key), string(val.Value))
		session.MarkMessage(val, "")
	}
	return nil
}
