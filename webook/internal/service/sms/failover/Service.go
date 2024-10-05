package failover

import (
	"context"
	"errors"
	"example.com/mod/webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{svcs: svcs}
}

// Send 失败了使用轮训
func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {

	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		//说明这个服务商出现了问题
		//可以打印日志
	}
	return errors.New("短信服务商全部失败")
}
