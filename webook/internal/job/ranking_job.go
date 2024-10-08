package job

import (
	"context"
	"example.com/mod/webook/internal/service"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func NewRankingJob(svc service.RankingService) Job {
	return &RankingJob{timeout: 2 * time.Second, svc: svc}
}

func (r *RankingJob) Name() string {

	return "RankingJob"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}
