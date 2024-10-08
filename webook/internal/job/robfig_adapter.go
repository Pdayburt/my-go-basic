package job

import (
	"go.uber.org/zap"
	"time"
)

type RankingJobAdapter struct {
	j Job
}

func NewRankingJobAdapter(job Job) *RankingJobAdapter {
	return &RankingJobAdapter{j: job}
}

func (r *RankingJobAdapter) Run() {
	zap.L().Info("Job任务开始执行～～",
		zap.String("name", r.j.Name()),
		zap.Any("time", time.Now()))
	err := r.j.Run()
	if err != nil {
		zap.L().Error("Job任务运行失败",
			zap.Error(err),
			zap.String("Job任务名称：", r.j.Name()))
	}
	zap.L().Info("Job任务执行成功～～",
		zap.String("name", r.j.Name()),
		zap.Any("time", time.Now()))
}
