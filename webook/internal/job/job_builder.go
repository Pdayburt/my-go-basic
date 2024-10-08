package job

/*import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronJobBuilder struct {
}

func (c *CronJobBuilder) Build(job Job) cron.Job {

	return cronJobFuncAdapter(func() error {
		name := job.Name()
		zap.L().Info("定时任务开始执行",
			zap.String("name", name))
		err := job.Run()
		if err != nil {
			zap.L().Info("定时任务执行失败",
				zap.String("name", name))
		}
		zap.L().Info("定时任务执行结束",
			zap.String("name", name))
		return nil
	})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}
*/
