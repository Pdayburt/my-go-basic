package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaSliceWindow string

// RedisSliceWindowLimiter redis上的额滑动窗口算法限流实现
type RedisSliceWindowLimiter struct {
	cmd redis.Cmdable
	//窗口大小
	interval time.Duration
	//阈值
	rate int
	//interval内允许rate个请求
}

func NewRedisSliceWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSliceWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSliceWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaSliceWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
