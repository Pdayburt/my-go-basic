package ioc

import (
	"example.com/mod/webook/internal/job"
	"example.com/mod/webook/internal/service"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService) job.Job {
	return job.NewRankingJob(svc)
}

func InitRankingJobAdapter(j job.Job) *job.RankingJobAdapter {
	return job.NewRankingJobAdapter(j)
}

func InitJobs(j *job.RankingJobAdapter) *cron.Cron {
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 4s", j)
	if err != nil {
		panic(err)
	}
	return expr

}
