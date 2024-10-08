package cache

import (
	"context"
	_ "embed"
	"example.com/mod/webook/interactive/domain"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

//go:embed lua/interative_incr_cnt.lua
var luaInteractive string

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	DecrReadCntIfPresent(ctx context.Context, biz string, id int64) error
	IncrCollectionCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, id int64, intr domain.Interactive) error
}

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, id int64, intr domain.Interactive) error {
	key := r.key(biz, id)
	err := r.client.HMSet(ctx, key,
		fieldLikeCnt, intr.Liked,
		fieldCollectCnt, intr.Collected,
		fieldReadCnt, intr.ReadCnt).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, 15*time.Minute).Err()
}

func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {

	data, err := r.client.HGetAll(ctx, r.key(biz, id)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(data) == 0 {
		//缓存不存在 系统错误
		return domain.Interactive{}, ErrKeyNotExists
	}
	collectCnt, _ := strconv.ParseInt(data[fieldCollectCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)
	return domain.Interactive{
		// 懒惰的写法
		CollectCnt: collectCnt,
		LikeCnt:    likeCnt,
		ReadCnt:    readCnt,
	}, err

}

func (r *RedisInteractiveCache) IncrCollectionCntIfPresent(ctx context.Context, biz string, id int64) error {
	return r.client.Eval(ctx, luaInteractive, []string{r.key(biz, id)},
		fieldCollectCnt, 1).Err()
}

func (r *RedisInteractiveCache) DecrReadCntIfPresent(ctx context.Context, biz string, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaInteractive,
		[]string{r.key(biz, bizId)},
		//read_cnt +1
		"read_cnt", 1).Err()
}

func NewInteractiveCache(redis redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{client: redis}
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
