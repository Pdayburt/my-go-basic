package web

import (
	"github.com/gin-gonic/gin"
)

type OAuth2WeChatHandler struct {
}

func NewOAuth2WeChatHandler() {

}

func (o *OAuth2WeChatHandler) ServeHTTP(gin *gin.Engine) {

	group := gin.Group("/oauth2/wechat")
	group.GET("/authurl", o.AuthUrl)
	group.GET("/callback", o.CallBack)

}

func (o *OAuth2WeChatHandler) AuthUrl(context *gin.Context) {

}

func (o *OAuth2WeChatHandler) CallBack(context *gin.Context) {

}
