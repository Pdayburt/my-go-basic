package time

import (
	"context"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	done := false
	for !done {
		select {
		case now := <-ticker.C:
			t.Log(now)
		case <-ctx.Done():
			t.Log("时间到了")
			done = true
		}
	}

}
