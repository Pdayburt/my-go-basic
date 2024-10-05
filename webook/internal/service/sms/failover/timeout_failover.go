package failover

import (
	"context"
	"example.com/mod/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailOverSMSService struct {
	svcs      []sms.Service
	idx       int32
	cnt       int32
	threshold int32
}

func (t *TimeoutFailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		nexIndex := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, nexIndex) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = nexIndex
	}
	service := t.svcs[idx]
	err := service.Send(ctx, tplId, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	default:
		return err
	}
	return err
}

func NewTimeoutFailOverSMSService() sms.Service {

	return &TimeoutFailOverSMSService{}
}
