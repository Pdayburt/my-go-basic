package failover

import (
	"context"
	"errors"
	"example.com/mod/webook/internal/service/sms"
	"example.com/mod/webook/internal/service/sms/ratelimitmocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailOverSMSService_Send(t *testing.T) {

	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		{
			name: "一次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {

				mockService0 := ratelimitmocks.NewMockService(ctrl)
				mockService0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{mockService0}

			},
			wantErr: nil,
		},
		{
			name: "重试成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {

				mockService0 := ratelimitmocks.NewMockService(ctrl)
				mockService0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				mockService1 := ratelimitmocks.NewMockService(ctrl)
				mockService1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{mockService0, mockService1}

			},
			wantErr: nil,
		},
		{
			name: "全部失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {

				mockService0 := ratelimitmocks.NewMockService(ctrl)
				mockService0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				mockService1 := ratelimitmocks.NewMockService(ctrl)
				mockService1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("发送全部失败"))
				return []sms.Service{mockService0, mockService1}

			},
			wantErr: errors.New("短信服务商全部失败"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFailOverSMSService := NewFailOverSMSService(tc.mock(ctrl))
			//ctx context.Context, tplId string, args []string, numbers ...string
			err := mockFailOverSMSService.Send(context.Background(), "mttpl",
				[]string{"1414"}, "562434")
			assert.Equal(t, tc.wantErr, err)

		})
	}
}
