package auth

import (
	"context"
	"example.com/mod/webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
}

func (s *SMSService) Send(ctx context.Context, tplId string,
	args []string, numbers ...string) error {
	//装饰器模式
	return s.svc.Send(ctx, tplId, args, numbers...)
}
