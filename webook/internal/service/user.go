package service

import (
	"context"
	"errors"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("邮箱或者账号错误")
	errInvalidUserOrPassword = errors.New("邮箱或者账号错误")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, u domain.User) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type UserServiceByRedisAndGORM struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceByRedisAndGORM{
		repo: repo,
	}
}

func (svc *UserServiceByRedisAndGORM) SignUp(ctx context.Context, u domain.User) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.repo.Create(ctx, u)
}

func (svc *UserServiceByRedisAndGORM) Login(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, u.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *UserServiceByRedisAndGORM) Profile(ctx context.Context, id int64) (domain.User, error) {
	byId, err := svc.repo.FindById(ctx, id)
	return byId, err
}

func (svc *UserServiceByRedisAndGORM) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return user, err
	}
	//没有该用户
	u := domain.User{Phone: phone}

	err = svc.repo.Create(ctx, u)
	if err != nil {
		return user, err
	}
	return svc.repo.FindByPhone(ctx, phone)

}
