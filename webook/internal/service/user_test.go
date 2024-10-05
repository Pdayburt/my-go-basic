package service

import (
	"context"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/repository"
	"example.com/mod/webook/internal/repository/repomocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
	"time"
)

func TestUserServiceByRedisAndGORM_Login(t *testing.T) {

	now := time.Now()
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		user     domain.User
		wantUser domain.User
		wantErr  error
	}{
		//FindByEmail(ctx context.Context, email string) (domain.User, error)
		{
			name: "测试-登陆成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				mockUserRepository := repomocks.NewMockUserRepository(ctrl)
				mockUserRepository.EXPECT().FindByEmail(gomock.Any(), "1123@qq.com").
					Return(domain.User{
						Email:    "1123@qq.com",
						Password: "$2a$10$ej8YvlNCB/tStJS9AkVcX.E8TP.iZ1JOhoxFZyFDYOwNYDaYSIg.W",
						Phone:    "12414134",
						Ctime:    now,
					}, nil)
				return mockUserRepository
			},
			ctx: context.Background(),
			user: domain.User{
				Email:    "1123@qq.com",
				Password: "123456#qqq",
			},
			wantUser: domain.User{
				Email:    "1123@qq.com",
				Password: "$2a$10$ej8YvlNCB/tStJS9AkVcX.E8TP.iZ1JOhoxFZyFDYOwNYDaYSIg.W",
				Phone:    "12414134",
				Ctime:    now,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCase {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := NewUserService(tc.mock(ctrl))
			user, err := userService.Login(tc.ctx, tc.user)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})

	}
}

func TestEncrypto(t *testing.T) {
	password, err := bcrypt.GenerateFromPassword([]byte("123456#qqq"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(string(password))
}
