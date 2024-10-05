package cache

import (
	"context"
	"encoding/json"
	"errors"
	"example.com/mod/webook/internal/domain"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrKeyNotExists = redis.Nil
)

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	client redis.Cmdable
	expiry time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client: client,
		expiry: time.Minute * 15,
	}
}

// Get 可以约定
// 只要err == nil 就认为缓存里有数据
// 如果没有数据 返回一个特定的error
func (uc *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := uc.key(id)
	val, err := uc.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err

}

func (uc *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := uc.key(user.Id)
	return uc.client.Set(ctx, key, bytes, uc.expiry).Err()
}

func (uc *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:id:%d", id)
}
