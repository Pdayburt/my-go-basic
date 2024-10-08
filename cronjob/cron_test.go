package cronjob

import (
	"github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestCron(t *testing.T) {

	expr := cron.New(cron.WithSeconds())
	expr.AddFunc("@every 3s", func() {
		t.Log("长任务开始", time.Now())
		time.Sleep(8 * time.Second)
		t.Log("长任务结束", time.Now())
	})
	//expr.AddJob("@every 3s", &MyJob{})
	expr.Start()
	time.Sleep(6 * time.Second)
	stop := expr.Stop()
	t.Log("停止信号已发出")

	<-stop.Done()
	t.Log("程序真正的结束")
}

type MyJob struct {
}

func (m *MyJob) Run() {
	log.Println("MyJob 运行了～～～")
}
