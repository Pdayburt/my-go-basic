package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type BatchConsumerHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
}

func NewBatchHandler[T any](fn func(msgs []*sarama.ConsumerMessage, t []T) error) *BatchConsumerHandler[T] {
	return &BatchConsumerHandler[T]{
		fn: fn,
	}
}

func (b *BatchConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	// 这个可以做成参数
	const batchSize = 10
	for {
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 这一批次已经超时了，
				// 或者，整个 consumer 被关闭了
				// 不再尝试凑够一批了
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					// channel 被关闭了
					return nil
				}
				msgs = append(msgs, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					// 消息格式都不对，没啥好处理的
					// 但是也不能直接返回，在线上的时候要继续处理下去
					zap.L().Error("反序列化消息体失败",
						zap.String("topic", msg.Topic),
						zap.Int32("partition", msg.Partition),
						zap.Int64("offset", msg.Offset),
						// 这里也可以考虑打印 msg.Value，但是有些时候 msg 本身也包含敏感数据
						zap.Error(err))
					// 不中断，继续下一个
					session.MarkMessage(msg, "")
					continue
				}
				ts = append(ts, t)
			}
		}
		err := b.fn(msgs, ts)
		if err == nil {
			// 这边就要都提交了
			for _, msg := range msgs {
				session.MarkMessage(msg, "")
			}
		} else {
			// 这里可以考虑重试，也可以在具体的业务逻辑里面重试
			// 也就是 eg.Go 里面重试
		}
		cancel()
	}

}
