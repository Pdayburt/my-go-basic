package web

import (
	"bytes"
	"errors"
	"example.com/mod/webook/internal/service"
	"example.com/mod/webook/internal/service/svcmocks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_RegisterRoutes(t *testing.T) {
	fmt.Println("TestUserHandler_RegisterRoutes")
}

func TestUserHandler_SignUp(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册功能-成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmocks.NewMockUserService(ctrl)
				mockUserService.EXPECT().SignUp(gomock.Any(), gomock.Any()).
					Return(nil)
				return mockUserService
			},
			reqBody: `
{
    "email": "1123@qq.com",
    "password": "123456#qqq",
    "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "注册功能-参数不对,bind",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmocks.NewMockUserService(ctrl)
				return mockUserService
			},
			reqBody: `
{
    "email": "1123@qq.com",
    "password": "123456#qqq",,
 "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "注册功能-邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmocks.NewMockUserService(ctrl)
				/*mockUserService.EXPECT().SignUp(gomock.Any(), gomock.Any()).
				Return(nil)*/
				return mockUserService
			},
			reqBody: `
{
    "email": "1123qq.com",
    "password": "123456#qqq",
    "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮件格式错误",
		},
		{
			name: "注册功能-两次密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmocks.NewMockUserService(ctrl)
				/*mockUserService.EXPECT().SignUp(gomock.Any(), gomock.Any()).
				Return(nil)*/
				return mockUserService
			},
			reqBody: `
{
    "email": "1123qq.com",
    "password": "12345146#qqq",
    "confirmPassword": "123456#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮件格式错误",
		},
		{
			name: "注册功能-邮箱已被注册",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmocks.NewMockUserService(ctrl)
				mockUserService.EXPECT().SignUp(gomock.Any(), gomock.Any()).
					Return(service.ErrUserDuplicateEmail)
				return mockUserService
			},
			reqBody: `
{
    "email": "1123@qq.com",
    "password": "12345146#qqq",
    "confirmPassword": "12345146#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱已被注册",
		},
		{
			name: "注册功能-系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				mockUserService := svcmocks.NewMockUserService(ctrl)
				mockUserService.EXPECT().SignUp(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
				return mockUserService
			},
			reqBody: `
{
    "email": "1123@qq.com",
    "password": "12345146#qqq",
    "confirmPassword": "12345146#qqq"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常！！",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userHandler := NewUserHandler(tc.mock(ctrl), nil)
			server := gin.Default()
			userHandler.RegisterRoutes(server)

			req, _ := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}

}

func TestMock(t *testing.T) {
}
