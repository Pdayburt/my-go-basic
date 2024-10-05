package repository

import (
	"context"
	"database/sql"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/repository/cache"
	"example.com/mod/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (ur *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.DomainToEntity(u))
}

func (ur *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.cache.Get(ctx, id)
	//缓存里有数据直接返回
	if err == nil {
		return u, nil
	}

	/*if errors.Is(err, cache.ErrKeyNotExists) {

	}*/
	//缓存里没这个数据 去数据库加载
	user, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = ur.EntityToDomain(user)
	//写入缓存
	err = ur.cache.Set(ctx, u)
	if err != nil {
		//打日志即可
		//缓存失败问题不大
	}
	return u, nil
}

func (ur *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	byEmail, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.EntityToDomain(byEmail), nil
}
func (ur *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	byEmail, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.EntityToDomain(byEmail), nil
}

func (ur *CachedUserRepository) EntityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Ctime:    time.UnixMilli(u.Ctime),
		Phone:    u.Phone.String,
	}
}
func (ur *CachedUserRepository) DomainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
}
