package web

import (
	"errors"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/service"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc      service.ArticleService
	interSvc service.InteractiveService
	biz      string
}

func NewArticleHandler(svc service.ArticleService, interSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		interSvc: interSvc,
		biz:      "article",
	}
}

func (ah *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/articles")
	group.POST("/edit", ah.Edit)
	group.POST("/publish", ah.Publish)
	group.POST("/withdraw", ah.withdraw)
	group.POST("/list", ah.list)
	group.GET("/detail/:id", ah.Detail)

	pub := server.Group("/pub")
	pub.GET("/:id", ah.PubDetail)

	////点赞和取消点赞
	pub.POST("/like", ah.Like)
}

func (ah *ArticleHandler) Like(ctx *gin.Context) {

	ctx, claim := ah.InitContextByUser(ctx)

	var likeReq LikeReq
	if err := ctx.Bind(&likeReq); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
			Data: fmt.Errorf("参数无法绑定%w", err),
		})
		return
	}
	var err error
	if likeReq.Like {
		err = ah.interSvc.Like(ctx, ah.biz, likeReq.Id, claim.Uid)
	} else {
		err = ah.interSvc.CancelLike(ctx, ah.biz, likeReq.Id, claim.Uid)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "Ok",
		Data: nil,
	})

}

func (ah *ArticleHandler) PubDetail(ctx *gin.Context) {

	ctx, claim := ah.InitContextByUser(ctx)
	idStr := ctx.Param("id")
	artId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("前端输入的 ID 不对", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
			Data: fmt.Errorf("查询文章详情的 ID %s 不正确, %w", idStr, err),
		})
		return
	}

	var eg errgroup.Group

	var art domain.Article
	eg.Go(func() error {
		art, err = ah.svc.GetPublishedById(ctx, claim.Uid, artId)
		return err
	})

	var intr domain.Interactive
	eg.Go(func() error {
		//如果这里的错误可以容能的话,记录日志即可 返回nil
		intr, err = ah.interSvc.Get(ctx, ah.biz, artId, claim.Uid)
		if err != nil {
			zap.L().Error("查询附属信息错误", zap.Error(err))
		}
		return nil
	})

	//这里会等待前两个完成
	//这里的err时上面两个中最先报错的err
	err = eg.Wait()
	if err != nil {
		//查询出错了
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("获取点赞等信息失败", zap.Error(err))

		return
	}
	//以下是异步执行
	go func() {
		er := ah.interSvc.IncrReadCnt(ctx, ah.biz, art.Id)
		if er != nil {
			zap.L().Error("增加阅读计数失败", zap.Int64("aid", art.Id), zap.Error(er))
		}
	}()

	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "查询成功",
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Abstract:   art.Abstract(),
			Content:    art.Content,
			Status:     art.Status.ToUint8(),
			Author:     art.Author.Name,
			Ctime:      art.Ctime.Format(time.RFC3339),
			Utime:      art.Utime.Format(time.RFC3339),
			Liked:      intr.Liked,
			Collected:  intr.Collected,
			LikeCnt:    intr.LikeCnt,
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
		},
	})

}

func (ah *ArticleHandler) Detail(ctx *gin.Context) {

	ctx, claim := ah.InitContextByUser(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		zap.L().Error("前端输入的 ID 不对", zap.Error(err))
		return
	}
	byId, err := ah.svc.GetById(ctx, claim.Uid, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "查询成功～～",
		Data: ArticleVo{
			Id:       byId.Id,
			Title:    byId.Title,
			Abstract: byId.Abstract(),
			Content:  byId.Content,
			Status:   byId.Status.ToUint8(),
			Author:   byId.Author.Name,
			Ctime:    byId.Ctime.Format(time.DateTime),
			Utime:    byId.Utime.Format(time.DateTime),
		},
	})

}

func (ah *ArticleHandler) list(ctx *gin.Context) {

	ctx, claim := ah.InitContextByUser(ctx)

	var req ListReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	res, err := ah.svc.List(ctx, claim.Uid, req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "获取成功",
		Data: slice.Map[domain.Article, ArticleVo](res,
			func(idx int, src domain.Article) ArticleVo {
				return ArticleVo{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					Ctime:    src.Ctime.Format(time.DateTime),
					Utime:    src.Utime.Format(time.DateTime),
				}
			}),
	})

}

func (ah *ArticleHandler) withdraw(ctx *gin.Context) {

	ctx, claim := ah.InitContextByUser(ctx)
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

	ctx, claim := ah.InitContextByUser(ctx)

	var req ArticleReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	id, err := ah.svc.Publish(ctx, req.ReqToDomain(claim.Uid))
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

	ctx, claim := ah.InitContextByUser(ctx)

	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	id, err := ah.svc.Save(ctx, req.ReqToDomain(claim.Uid))
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

func (ah *ArticleHandler) InitContextByUser(ctx *gin.Context) (*gin.Context, UserClaims) {
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
		return nil, UserClaims{}
	}
	return ctx, claim

}
