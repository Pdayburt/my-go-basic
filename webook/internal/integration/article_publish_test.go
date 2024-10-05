package integration

import (
	"bytes"
	"encoding/json"
	"example.com/mod/webook/internal/integration/startup"
	"example.com/mod/webook/internal/ioc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticlePublishTestSuit struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (ap *ArticlePublishTestSuit) SetupSuite() {
	ap.server = gin.Default()
	articleHandler := startup.InitArticleHandler()
	articleHandler.RegisterRoutes(ap.server)
	ap.db = ioc.InitDbB()
}

func (ap *ArticlePublishTestSuit) TestPublish() {
	t := ap.T()
	testCase := []struct {
		name string
		//集成测试准备数据
		before func(t *testing.T)
		////集成测试验证数据
		after func(t *testing.T)
		art   Article
		//http响应码
		wantCode int
		//我希望http响应带上帖子的id
		wantRes Result[int64]
	}{
		{
			name: "新建文章并发表",
			before: func(t *testing.T) {
				ap.server.Use(func(ctx *gin.Context) {

				})
			},
			art: Article{
				Title:   "我的标题-新建文章并发表",
				Content: "我的内容-新建文章并发表",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 1,
				Msg:  "OK",
				Data: 266,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			resp := httptest.NewRecorder()
			ap.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				t.Log("response code:", resp.Code)
				return
			}
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}

}

func TestArticlePublish(t *testing.T) {
	suite.Run(t, &ArticlePublishTestSuit{})
}
