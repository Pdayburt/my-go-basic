package article

import (
	"context"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewInteractiveReadEventConsumer(cl sarama.Client, repo repository.InteractiveRepository) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: cl,
		repo:   repo,
	}

}

func (k *InteractiveReadEventConsumer) Start() error {

	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"article_read"},
			saramax.NewHandler[ReadEvent](k.Consume))
		if err != nil {
			zap.L().Error("退出消费循环一场", zap.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	return k.repo.IncrReadCnt(ctx, "article", t.Aid)
}
