package article

import (
	"context"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewInteractiveReadEventBatchConsumer(cl sarama.Client, repo repository.InteractiveRepository) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: cl,
		repo:   repo,
	}
}

func (k *InteractiveReadEventBatchConsumer) Start() error {

	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}
	go func() {
		//ctx context.Context, topics []string, handler ConsumerGroupHandler
		err := cg.Consume(context.Background(),
			[]string{"article_read"},
			saramax.NewBatchHandler[ReadEvent](k.Consume))
		if err != nil {
			zap.L().Error("退出消费循环一场", zap.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventBatchConsumer) Consume(msg []*sarama.ConsumerMessage, t []ReadEvent) error {
	ids := make([]int64, 0, len(t))
	bizs := make([]string, 0, len(t))
	for _, evt := range t {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	err := k.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		zap.L().Error("批量插入失败",
			zap.Int64s("ids", ids),
			zap.Error(err))
	}
	return nil
}
