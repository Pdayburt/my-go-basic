package ratelimit

import (
	"context"
	"example.com/mod/webook/internal/service/sms"
	"example.com/mod/webook/pkg/ratelimit"
	"fmt"
)

// RateLimitSMSServiceV1
// 使用组合：
// • 用户可以直接访问 Service，绕开你装饰器 本身。
// • 可以只实现 Service 的部分方法。
// 不使用组合：
//• 可以有效阻止用户绕开装饰器。
//• 必须实现 Service 的全部方法。

type RateLimitSMSServiceV1 struct {
	sms.Service
	limit ratelimit.Limiter
}

func NewServiceV1(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitSMSServiceV1{
		Service: svc,
		limit:   limiter,
	}
}

func (s *RateLimitSMSServiceV1) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
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

	err = s.Service.Send(ctx, tplId, args, numbers...)

	//方法后➕新特性
	return err
}
