package cache

import (
	"context"
	"example.com/mod/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {

	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "验证码存储成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mockCmdable := redismocks.NewMockCmdable(ctrl)
				//i, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int64()
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))

				mockCmdable.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:15234124"}, "123456").Return(cmd)

				return mockCmdable
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "15234124",
			code:    "123456",
			wantErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctrl.Finish()
			mockCodeCache := NewCodeCache(tc.mock(ctrl))
			//ctx context.Context, biz, phone, code string
			err := mockCodeCache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)

		})
	}

}
