package ratelimit

import (
	"context"
	"example.com/mod/webook/internal/service/sms"
	"example.com/mod/webook/pkg/ratelimit"
	"fmt"
)

var (
	errLimit = fmt.Errorf("触发了限流")
)

type RateLimitSMSService struct {
	svc   sms.Service
	limit ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitSMSService{
		svc:   svc,
		limit: limiter,
	}
}

func (s *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//方法前➕新特性

	limit, err := s.limit.Limit(ctx, "sms:tencent")
	if err != nil {
		//系统错误
		fmt.Errorf("短信服务判断是否限流出现问题%w", err)
	}
	//	触发了限流
	if limit {
		return errLimit
	}

	err = s.svc.Send(ctx, tplId, args, numbers...)

	//方法后➕新特性
	return err
}
