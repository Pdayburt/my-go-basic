package integration

import (
	"bytes"
	"encoding/json"
	"example.com/mod/webook/internal/integration/startup"
	"example.com/mod/webook/internal/ioc"
	"example.com/mod/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 定义测试套件结构:
type ArticleSaveTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (at *ArticleSaveTestSuite) SetupSuite() {
	at.server = gin.Default()
	articleHandler := startup.InitArticleHandler()
	articleHandler.RegisterRoutes(at.server)
	at.db = ioc.InitDbB()
}

/*
func (at *ArticleSaveTestSuite) TearDownSuite() {
	//清空所有数据，并且自增主键恢复到1


}*/

func (at *ArticleSaveTestSuite) TestEdit() {

	t := at.T()
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
			name: "新建帖子-保存成功",
			before: func(t *testing.T) {

				/*at.server.Use(func(context *gin.Context) {
					context.Set("user", web.UserClaims{
						Uid: 266,
					})
				})*/
				//	at.db.Exec("TRUNCATE TABLE articles")
			},
			after: func(t *testing.T) {
				//验证数据库
				var art dao.Article
				err := at.db.Find(&art, "id = ?", 4).Error
				if err != nil {
					t.Fatal(err)
				}

				assert.True(t, art.Utime > 0)
				assert.True(t, art.Ctime > 0)
				art.Ctime = 0
				art.Utime = 0

				assert.Equal(t, dao.Article{
					Id:       4,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 266,
				}, art)
			},

			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 1,
				Msg:  "OK",
				Data: 266,
			},
		},
		{
			name: "修改已有帖子，并保存",
			before: func(t *testing.T) {

				at.db.Exec("TRUNCATE TABLE articles")
				err := at.db.Create(&dao.Article{
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 266,
					Ctime:    123,
					Utime:    456,
				}).Error
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				//验证数据库
				var art dao.Article
				err := at.db.Find(&art, "id = ?", 1).Error
				if err != nil {
					t.Fatal(err)
				}
				assert.True(t, art.Utime > 456)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的新的标题",
					Content:  "我的新的内容",
					AuthorId: 266,
					Ctime:    123,
					Utime:    0,
				}, art)
			},

			art: Article{
				Id:      1,
				Title:   "我的新的标题",
				Content: "我的新的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 1,
				Msg:  "OK",
				Data: 1,
			},
		},
		{
			name: "修改别人的帖子～～",
			before: func(t *testing.T) {

				at.db.Exec("TRUNCATE TABLE articles")
				err := at.db.Create(&dao.Article{
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 789,
					Ctime:    123,
					Utime:    456,
				}).Error
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				//验证数据库
				var art dao.Article
				err := at.db.Find(&art, "id = ?", 1).Error
				if err != nil {
					t.Fatal(err)
				}
				assert.True(t, art.Utime > 456)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "我的新的标题",
					Content:  "我的新的内容",
					AuthorId: 266,
					Ctime:    123,
					Utime:    0,
				}, art)
			},

			art: Article{
				Id:      1,
				Title:   "我的新的标题",
				Content: "我的新的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 1,
				Msg:  "OK",
				Data: 1,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			//构造请求
			//执行
			//验证
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/articles/edit",
				bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			at.server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				t.Log("response code:", resp.Code)
				return
			}

			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)

		})
	}

}

func TestArticle(t *testing.T) {

	suite.Run(t, &ArticleSaveTestSuite{})

}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
