package service

import (
	"context"
	intrv1 "example.com/mod/webook/api/proto/gen/intr/v1"
	"example.com/mod/webook/internal/domain"
	"github.com/ecodeclub/ekit/queue"
	"go.uber.org/zap"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	topN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	artSvc    ArticleService
	intrSvc   intrv1.InteractiveServiceClient
	batchSize int
	n         int
	scoreFunc func(t time.Time, likeCnt int64) float64
}

func NewBatchRankingService(artSvc ArticleService, intrSvc intrv1.InteractiveServiceClient) RankingService {
	return &BatchRankingService{artSvc: artSvc, intrSvc: intrSvc}
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	zap.L().Info("～～ranking任务执行中～～", zap.Any("time", time.Now()))
	return nil
}

func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	//拿一批数据
	offSet := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	_ = queue.NewPriorityQueue[Score](b.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score < dst.score {
				return -1
			} else {
				return 0
			}
		})

	arts, err := b.artSvc.ListPub(ctx, offSet, b.batchSize)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(arts))
	for _, art := range arts {
		ids = append(ids, art.Id)
	}

	_, err = b.intrSvc.GetByIds(ctx, &intrv1.GetByIdsReq{
		Biz: "article",
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	//计算score
	return arts, nil
}
