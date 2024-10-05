package web

import (
	"errors"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (ah *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/articles")
	group.POST("/edit", ah.Edit)
	group.POST("/publish", ah.Publish)
	group.POST("/withdraw", ah.withdraw)
}

func (ah *ArticleHandler) withdraw(ctx *gin.Context) {

	//测试用
	ctx.Set("user", UserClaims{
		Uid: 266,
	})
	c := ctx.MustGet("user")
	claim, ok := c.(UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("保存失败", zap.Error(errors.New("未发现用户的id")))
		return
	}

	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		return
	}
	err := ah.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claim.Uid,
		},
	})

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("发表贴失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (ah *ArticleHandler) Publish(ctx *gin.Context) {

	var req ArticleReq

	if err := ctx.Bind(&req); err != nil {
		return
	}
	//测试用
	ctx.Set("user", UserClaims{
		Uid: 266,
	})

	c := ctx.MustGet("user")
	claims, ok := c.(UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("保存失败", zap.Error(errors.New("未发现用户的id")))
		return
	}

	id, err := ah.svc.Publish(ctx, req.ReqToDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("发表贴失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 1,
		Msg:  "OK",
		Data: id,
	})

}

func (ah *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("user")
	claims, ok := c.(UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("保存失败", zap.Error(errors.New("未发现用户的id")))
		return
	}

	id, err := ah.svc.Save(ctx, req.ReqToDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("保存失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 1,
		Msg:  "OK",
		Data: id,
	})
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (ar ArticleReq) ReqToDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      ar.Id,
		Title:   ar.Title,
		Content: ar.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
