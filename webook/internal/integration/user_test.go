package integration

import (
	"bytes"
	"encoding/json"
	"example.com/mod/webook/internal/integration/startup"
	"example.com/mod/webook/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {

	ginServer := startup.InitWebServerByWire()

	testCase := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//不需要，也就是redis里什么数据都没有
			},
			after: func(t *testing.T) {

			},
			reqBody: `{
    "phone": "009788660"
}`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 0,
				Msg:  "发送成功",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			ginServer.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			var res web.Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.Equal(t, tc.wantBody, res)

		})
	}

}
