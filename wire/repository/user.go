package repository

import "example.com/mod/wire/repository/dao"

type UserRepository struct {
	sd *dao.UserDao
}

func NewUserRepository(sd *dao.UserDao) *UserRepository {
	return &UserRepository{sd: sd}
}
